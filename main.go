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

package main

import (
	"os"
	"os/signal"
	"strconv"

	"github.com/flexpool/solo/configuration"
	"github.com/flexpool/solo/engine"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/process"
	"github.com/flexpool/solo/utils"
	"github.com/sirupsen/logrus"
)

func interruptHandler(engine *engine.MiningEngine) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for i := 0; i < 10; i++ {
		<-interrupt
		if i == 0 {
			log.Logger.Info("Caught interrupt, exitting...")
			process.ExitChan <- 2
		} else {
			log.Logger.Info("Caught interrupt. Press Ctrl+C " + strconv.Itoa(10-i) + " times to panic.")
		}
	}

	panic("Force shutdown by user")
}

func main() {
	// Initialize logging
	log.InitLog()

	// Get the configuration
	config, err := configuration.GetConfig()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "config",
			"error":  err,
		}).Fatal("Unable to get config")
	}

	// Set log level
	log.SetLogLevel(config.LogLevel)

	// Check the config
	err = utils.IsInvalidAddress(config.WorkmanagerNotificationsBindAddr)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "config",
			"error":  err,
		}).Error("Invalid Work Receiver bind address")
		os.Exit(1)
	}

	if config.GatewayInsecureBindAddr == "" /* && config.GatewaySSLBindAddr == "" */ {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "config",
		}).Error("At least one gateway bind address should be specified")
		os.Exit(1)
	}

	miningEngine, err := engine.NewMiningEngine(config.WorkmanagerNotificationsBindAddr, config.ShareDifficulty, config.GatewayInsecureBindAddr, "", config.GatewayPassword, config.NodeHTTPRPC, config.DBPath)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "engine",
			"error":  err.Error(),
		}).Error("Unable to initialize the mining engine")
		os.Exit(1)
	}

	miningEngine.Start()

	go interruptHandler(miningEngine)
	exitCode := <-process.ExitChan

	miningEngine.Stop()
	log.Logger.WithFields(logrus.Fields{
		"prefix": "engine",
	}).Info("Stopped mining engine")

	os.Exit(exitCode)
}
