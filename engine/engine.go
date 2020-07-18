package engine

import (
	"github.com/flexpool/solo/gateway"
	"github.com/flexpool/solo/log"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// MiningEngine represents the Flexpool Solo mining engine
type MiningEngine struct {
	Workreceiver     *gateway.WorkReceiver
	workreceiverBind string
	shareDifficulty  uint64
	Gateways         []*gateway.Gateway
}

// NewMiningEngine creates a new Mining Engine
func NewMiningEngine(workreceiverBind string, shareDifficulty uint64, insecureStratumBind string, secureStratumBind string, stratumPassword string) (*MiningEngine, error) {
	engine := MiningEngine{
		Workreceiver:     gateway.NewWorkReceiver(workreceiverBind, shareDifficulty),
		workreceiverBind: workreceiverBind,
		shareDifficulty:  shareDifficulty,
	}

	if insecureStratumBind != "" {
		gatewayInsecure, err := gateway.NewGatewayInsecure(engine.Workreceiver, insecureStratumBind, stratumPassword)
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
	go e.Workreceiver.Run()

	log.Logger.WithFields(logrus.Fields{
		"prefix": "engine",
		"bind":   e.workreceiverBind,
	}).Info("Started Work Receiver server")

	for _, g := range e.Gateways {
		go g.Run()
	}
}

// Stop stops the mining engine
func (e *MiningEngine) Stop() {
	e.Workreceiver.Stop()
	for _, g := range e.Gateways {
		g.Stop()
	}
}
