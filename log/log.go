package log

import (
	"strings"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// Logger is a logrus logging instance for the Solo application
var Logger = logrus.New()

// InitLog initializes Logger
func InitLog() {
	formater := new(prefixed.TextFormatter)
	// Set timestamp format to YYYY-MM-DD HH:MM:SS
	formater.TimestampFormat = "2006-01-02 15:04:05"
	formater.FullTimestamp = true
	Logger.Formatter = formater
}

// SetLogLevel sets the log level
func SetLogLevel(logLevelString string) {
	switch strings.ToLower(logLevelString) {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "info":
		Logger.SetLevel(logrus.DebugLevel)
	case "warning":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.Warn("Unknown log level \"" + logLevelString + "\". Falling back to INFO.")
	}
}
