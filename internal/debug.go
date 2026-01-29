package internal

import (
	"log"
	"os"
)

var (
	logFlags = log.Ldate | log.Ltime
	info     = log.New(os.Stdout, "INFO: ", logFlags)
	debug    = log.New(os.Stdout, "DEBUG: ", logFlags)
	errorLog = log.New(os.Stderr, "ERROR: ", logFlags)
	warning  = log.New(os.Stderr, "WARNING: ", logFlags)
	plain    = log.New(os.Stdout, "", 0)
)

// DebugLogger provides an interface to call the logging functions
type DebugLogger interface {
	LogInfo(f string, v ...any)
	LogDebug(f string, v ...any)
	LogError(f string, v ...any)
	LogWarning(f string, v ...any)
	PrintVersion(f string, v ...any)
}

// Debug contains information about the debug level
type Debug struct {
	DebugLevel bool
}

// LogInfo prints info-level messages to standard out
func (d Debug) LogInfo(f string, v ...any) {
	info.Printf(f, v...)
}

// LogDebug prints debug-level messages to standard out
func (d Debug) LogDebug(f string, v ...interface{}) {
	if d.DebugLevel {
		debug.Printf(f, v...)
	}
}

// LogError prints error-level messages to standard out
func (d Debug) LogError(f string, v ...interface{}) {
	errorLog.Printf(f, v...)
}

// LogWarning prints warning-level messages to standard out
func (d Debug) LogWarning(f string, v ...interface{}) {
	if d.DebugLevel {
		warning.Printf(f, v...)
	}
}

// PrintVersion prints the TeleIRC version number
func (d Debug) PrintVersion(f string, v ...interface{}) {
	plain.Printf(f, v...)
}
