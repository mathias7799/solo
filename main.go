package main

import (
	"os"
	"os/signal"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/flexpool/solo/configuration"
	"github.com/flexpool/solo/engine"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/process"
	"github.com/flexpool/solo/utils"
	"github.com/sirupsen/logrus"
)

func interruptHandler() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for i := 0; i < 10; i++ {
		<-interrupt
		if i == 0 {
			log.Logger.Info("Caught interrupt, exitting...")
			process.ExitChan <- 2
		} else {
			log.Logger.Info("Caught interrupt. Press Ctrl+C " + strconv.Itoa(i) + " times to panic.")
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

	// Check the config
	err = utils.IsInvalidAddress(config.WorkreceiverBindAddr)
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

	miningEngine, err := engine.NewMiningEngine(config.WorkreceiverBindAddr, config.ShareDifficulty, config.GatewayInsecureBindAddr, "", config.GatewayPassword)

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "engine",
			"error":  err,
		}).Error("Unable to create a mining engine")
	}

	miningEngine.Start()
	log.Logger.WithFields(logrus.Fields{
		"prefix":     "engine",
		"share-diff": humanize.SIWithDigits(float64(config.ShareDifficulty), 2, "H"),
	}).Info("Started mining engine")

	go interruptHandler()
	exitCode := <-process.ExitChan
	miningEngine.Stop()
	log.Logger.WithFields(logrus.Fields{
		"prefix": "engine",
	}).Info("Stopped mining engine")
	os.Exit(exitCode)
}
