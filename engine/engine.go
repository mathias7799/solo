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

package engine

import (
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/flexpool/solo/db"
	"github.com/flexpool/solo/gateway"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/nodeapi"
	"github.com/flexpool/solo/stats"
	"github.com/flexpool/solo/web"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// MiningEngine represents the Flexpool Solo mining engine
type MiningEngine struct {
	workmanagerNotificationsBind string
	webServerBind                string
	shareDifficulty              uint64

	Workmanager              *gateway.WorkManager
	Gateways                 []*gateway.Gateway
	StatsCollector           *stats.Collector
	Database                 *db.Database
	BlockConfirmationManager *stats.BlockConfirmationManager

	WebServer *web.Server

	waitGroup *sync.WaitGroup
}

// NewMiningEngine creates a new Mining Engine
func NewMiningEngine(workmanagerNotificationsBind string, shareDifficulty uint64, insecureStratumBind string, secureStratumBind string, stratumPassword string, nodeHTTPRPC string, databasePath string, blockConfirmationsRequired uint64, webServerBind string) (*MiningEngine, error) {
	node, err := nodeapi.NewNode(nodeHTTPRPC)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Node")
	}

	database, err := db.OpenDB(databasePath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open db")
	}

	waitGroup := new(sync.WaitGroup)

	statsCollector := stats.NewCollector(database, waitGroup, shareDifficulty)
	blockConfirmationManager := stats.NewBlockConfirmationManager(database, waitGroup, node, blockConfirmationsRequired)

	workmanager := gateway.NewWorkManager(workmanagerNotificationsBind, shareDifficulty, node, waitGroup)

	webServer := web.NewServer(database, node, waitGroup, workmanager, webServerBind)

	engine := MiningEngine{
		Workmanager:                  workmanager,
		workmanagerNotificationsBind: workmanagerNotificationsBind,
		shareDifficulty:              shareDifficulty,
		StatsCollector:               statsCollector,
		BlockConfirmationManager:     blockConfirmationManager,
		Database:                     database,
		waitGroup:                    waitGroup,
		WebServer:                    webServer,
		webServerBind:                webServerBind,
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

	go e.BlockConfirmationManager.Run()

	go e.WebServer.Run()
	log.Logger.WithFields(logrus.Fields{
		"prefix": "engine",
		"bind":   e.webServerBind,
	}).Info("Started Web Server")

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
	e.BlockConfirmationManager.Stop()
	e.WebServer.Stop()

	e.waitGroup.Wait()
	e.Database.DB.Close()

	log.Logger.WithFields(logrus.Fields{
		"prefix": "engine",
	}).Info("Closed database")
}
