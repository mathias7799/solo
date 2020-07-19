// Flexpool Solo - A lightweight SOLO Ethereum mining pool
// Copyright (C) 2020  Flexpool
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
