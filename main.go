package main

import (
	"os"
	"os/signal"

	"github.com/flexpool/solo/configuration"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/workreceiver"
	"github.com/sirupsen/logrus"
)

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

	workReceiver := workreceiver.NewWorkReceiver(config.WorkreceiverBindAddr)
	go workReceiver.Start()
	if workReceiver.Exited {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "workreceiver",
			"error":  err,
		}).Fatal("Unable to start Work Receiver")
	}

	log.Logger.WithFields(logrus.Fields{
		"prefix": "workreceiver",
		"bind":   config.WorkreceiverBindAddr,
	}).Info("Started Work Receiver server")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	log.Logger.Info("Caught interrupt, exitting...")
}
