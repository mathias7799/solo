package engine

import (
	"github.com/flexpool/solo/gateway"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/nodeapi"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// MiningEngine represents the Flexpool Solo mining engine
type MiningEngine struct {
	Workmanager                  *gateway.WorkManager
	workmanagerNotificationsBind string
	shareDifficulty              uint64
	Gateways                     []*gateway.Gateway
}

// NewMiningEngine creates a new Mining Engine
func NewMiningEngine(workreceiverBind string, shareDifficulty uint64, insecureStratumBind string, secureStratumBind string, stratumPassword string, nodeHTTPRPC string) (*MiningEngine, error) {
	node, err := nodeapi.NewNode(nodeHTTPRPC)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Node")
	}

	engine := MiningEngine{
		Workmanager:                  gateway.NewWorkManager(workreceiverBind, shareDifficulty, node),
		workmanagerNotificationsBind: workreceiverBind,
		shareDifficulty:              shareDifficulty,
	}

	if insecureStratumBind != "" {
		gatewayInsecure, err := gateway.NewGatewayInsecure(engine.Workmanager, insecureStratumBind, stratumPassword)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to initialize insecure gateway")
		}
		engine.Gateways = append(engine.Gateways, &gatewayInsecure)
	}

	if secureStratumBind != "" {
		return nil, errors.New("secure stratum is unimplemented")
	}

	return &engine, nil
}

// Start starts the mining engine
func (e *MiningEngine) Start() {
	// Starting work receiver
	go e.Workmanager.Run()

	log.Logger.WithFields(logrus.Fields{
		"prefix":             "engine",
		"notifications-bind": e.workmanagerNotificationsBind,
	}).Info("Started Work Manager")

	for _, g := range e.Gateways {
		go g.Run()
	}
}

// Stop stops the mining engine
func (e *MiningEngine) Stop() {
	e.Workmanager.Stop()
	for _, g := range e.Gateways {
		g.Stop()
	}
}
