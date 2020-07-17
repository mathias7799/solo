package log

import (
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
