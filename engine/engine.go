package engine

import (
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/flexpool/solo/db"
	"github.com/flexpool/solo/gateway"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/nodeapi"
	"github.com/flexpool/solo/stats"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// MiningEngine represents the Flexpool Solo mining engine
type MiningEngine struct {
	Workmanager                  *gateway.WorkManager
	workmanagerNotificationsBind string
	shareDifficulty              uint64
	Gateways                     []*gateway.Gateway
	StatsCollector               *stats.Collector
	Database                     *db.Database
	waitGroup                    *sync.WaitGroup
}

// NewMiningEngine creates a new Mining Engine
func NewMiningEngine(workmanagerNotificationsBind string, shareDifficulty uint64, insecureStratumBind string, secureStratumBind string, stratumPassword string, nodeHTTPRPC string, databasePath string) (*MiningEngine, error) {
	node, err := nodeapi.NewNode(nodeHTTPRPC)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Node")
	}

	database, err := db.OpenDB(databasePath)
	if err != nil {
		return nil, err
	}

	waitGroup := new(sync.WaitGroup)

	statsCollector := stats.NewCollector(database, waitGroup)

	engine := MiningEngine{
		Workmanager:                  gateway.NewWorkManager(workmanagerNotificationsBind, shareDifficulty, node, waitGroup),
		workmanagerNotificationsBind: workmanagerNotificationsBind,
		shareDifficulty:              shareDifficulty,
		StatsCollector:               statsCollector,
		Database:                     database,
		waitGroup:                    waitGroup,
	}

	if insecureStratumBind != "" {
		gatewayInsecure, err := gateway.NewGatewayInsecure(engine.Workmanager, insecureStratumBind, stratumPassword, engine.StatsCollector, waitGroup)
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
	// Starting work manager
	go e.Workmanager.Run()

	log.Logger.WithFields(logrus.Fields{
		"prefix":             "engine",
		"notifications-bind": e.workmanagerNotificationsBind,
	}).Info("Started Work Manager")

	go e.StatsCollector.Run()

	for _, g := range e.Gateways {
		go g.Run()
	}

	log.Logger.WithFields(logrus.Fields{
		"prefix":     "engine",
		"share-diff": humanize.SIWithDigits(float64(e.shareDifficulty), 2, "H"),
	}).Info("Started mining engine")
}

// Stop stops the mining engine
func (e *MiningEngine) Stop() {
	for _, g := range e.Gateways {
		g.Stop()
	}
	e.Workmanager.Stop()
	e.StatsCollector.Stop()

	e.waitGroup.Wait()
	e.Database.DB.Close()

	log.Logger.WithFields(logrus.Fields{
		"prefix": "engine",
	}).Info("Closed database")
}
