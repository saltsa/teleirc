// Package telegram handles all Telegram-side logic.
package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram/bot"
	"github.com/ritlug/teleirc/internal"
)

/*
Client contains information for the Telegram bridge, including
the TelegramSettings needed to run the bot
*/
type Client struct {
	API           *tgbotapi.Bot
	Settings      *internal.TelegramSettings
	IRCSettings   *internal.IRCSettings
	ImgurSettings *internal.ImgurSettings
	logger        internal.DebugLogger
	sendToIrc     func(string)

	ctx       context.Context
	ctxCancel context.CancelFunc
}

/*
NewClient creates a new Telegram bot client
*/
func NewClient(settings *internal.TelegramSettings, ircsettings *internal.IRCSettings, imgur *internal.ImgurSettings, logger internal.DebugLogger) *Client {
	logger.LogInfo("Creating new Telegram bot client...")
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{ctx: ctx, ctxCancel: cancel, Settings: settings, IRCSettings: ircsettings, ImgurSettings: imgur, logger: logger}
}

/*
SendMessage sends a message to the Telegram channel specified in the settings
*/
func (tg *Client) SendMessage(msg string) {
	tg.logger.LogDebug("tg send message: %s", msg)
	newMsg := &tgbotapi.SendMessageParams{
		ChatID: tg.Settings.ChatID,
		Text:   msg,
	}

	if _, err := tg.API.SendMessage(tg.ctx, newMsg); err != nil {
		var attempts int = 0
		// Try resending 3 times if the message is successfully sent
		for err != nil && attempts < 3 {
			attempts++
			tg.logger.LogError("send failure #%d: %s", err, attempts)
			_, err = tg.API.SendMessage(tg.ctx, newMsg)
		}
	}
}

/*
StartBot adds necessary handlers to the client and then connects,
returning any errors that occur
*/
func (tg *Client) StartBot(errChan chan<- error, sendMessage func(string)) {
	tg.logger.LogInfo("Starting up Telegram bot...")
	var err error

	opts := []tgbotapi.Option{
		tgbotapi.WithDefaultHandler(messageHandler(tg)),
		tgbotapi.WithSkipGetMe(),
	}
	if tg.Settings.DebugEnabled {
		opts = append(opts, tgbotapi.WithDebug())
	}

	tg.API, err = tgbotapi.New(tg.Settings.Token, opts...)
	if err != nil {
		errChan <- err
	}

	me, err := tg.API.GetMe(tg.ctx)
	if err != nil {
		errChan <- err
	}
	tg.logger.LogInfo("Authorized on account %s", me.Username)
	tg.sendToIrc = sendMessage

	tg.API.Start(tg.ctx)

	errChan <- nil
}

func (tg *Client) Close() {
	tg.ctxCancel()
}
