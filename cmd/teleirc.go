// Package main contains all logic relating to running TeleIRC
package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/ritlug/teleirc/internal"
	"github.com/ritlug/teleirc/internal/handlers/irc"
	tg "github.com/ritlug/teleirc/internal/handlers/telegram"
)

var (
	flagPath    = flag.String("conf", ".env", "config file")
	flagDebug   = flag.Bool("debug", false, "enable debugging output")
	flagVersion = flag.Bool("version", false, "displays current version of TeleIRC")
	version     string
)

func main() {
	flag.Parse()
	logger := internal.Debug{DebugLevel: *flagDebug}

	logger.LogInfo("Current TeleIRC version: %s", version)

	if *flagVersion {
		return
	}

	// Notify that logger is enabled
	logger.LogDebug("Debug mode enabled!")

	settings, err := internal.LoadConfig(*flagPath)
	if err != nil {
		logger.LogError("config load: %s", err)
		os.Exit(1)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	tgClient := tg.NewClient(&settings.Telegram, &settings.IRC, &settings.Imgur, logger)
	tgChan := make(chan error)

	ircClient := irc.NewClient(&settings.IRC, &settings.Telegram, logger)
	ircChan := make(chan error)

	go ircClient.StartBot(ircChan, tgClient.SendMessage)
	go tgClient.StartBot(tgChan, ircClient.SendMessage)

	exitError := false

	select {
	case ircErr := <-ircChan:
		logger.LogError("IRC error: %s", ircErr)
		exitError = true
	case tgErr := <-tgChan:
		logger.LogError("Telegram error: %s", tgErr)
		exitError = true
	case signal := <-signalChannel:
		logger.LogInfo("Signal Received: %s", signal.String())
		break
	}

	logger.LogInfo("Shutting Down...")
	ircClient.Close()
	tgClient.Close()
	logger.LogInfo("Exiting")

	if exitError {
		os.Exit(1)
	}
}
