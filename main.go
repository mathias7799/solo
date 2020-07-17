package main

import (
	"os"
	"os/signal"
	"strconv"

	"github.com/flexpool/solo/configuration"
	"github.com/flexpool/solo/gateway"
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

	workReceiver := gateway.NewWorkReceiver(config.WorkreceiverBindAddr, config.ShareDifficulty)
	go workReceiver.Run()

	log.Logger.WithFields(logrus.Fields{
		"prefix": "workreceiver",
		"bind":   config.WorkreceiverBindAddr,
	}).Info("Started Work Receiver server")

	var gateways []gateway.Gateway
	var gatewayTmp gateway.Gateway

	if config.GatewayInsecureBindAddr != "" {
		gatewayTmp, err = gateway.NewGatewayInsecure(workReceiver, config.GatewayInsecureBindAddr, config.GatewayPassword)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"prefix": "gateway",
				"bind":   config.GatewayInsecureBindAddr,
				"secure": "false",
			}).Error("Unable to start gateway")
			workReceiver.Stop()
			os.Exit(1)
		}
		gateways = append(gateways, gatewayTmp)
	}

	for _, gateway := range gateways {
		go gateway.Run()
	}

	log.Logger.WithFields(logrus.Fields{
		"prefix":     "workreceiver",
		"share-diff": config.ShareDifficulty,
	}).Info("Initialized mining engine")

	go interruptHandler()
	exitCode := <-process.ExitChan
	workReceiver.Stop()
	for _, gateway := range gateways {
		gateway.Stop()
	}
	os.Exit(exitCode)
}
