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

package gateway

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"sync"
	"time"

	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/process"
	"github.com/flexpool/solo/stats"
	"github.com/flexpool/solo/utils"
	"github.com/sirupsen/logrus"
)

// Gateway is a stratum proxy that servers workers
type Gateway struct {
	bind              string
	stratumPassword   string
	isSecure          bool
	tlsKeyPair        tls.Certificate
	context           context.Context
	cancelContextFunc context.CancelFunc
	parentWorkManager *WorkManager
	statsCollector    *stats.Collector
	engineWaitGroup   *sync.WaitGroup
}

// NewGatewayInsecure creates Non SSL gateway instance
func NewGatewayInsecure(parentWorkManager *WorkManager, bind string, password string, statsCollector *stats.Collector, engineWaitGroup *sync.WaitGroup) (Gateway, error) {
	err := utils.IsInvalidAddress(bind)
	if err != nil {
		return Gateway{}, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	return Gateway{bind: bind, stratumPassword: password, isSecure: false, context: ctx, cancelContextFunc: cancelFunc, parentWorkManager: parentWorkManager, statsCollector: statsCollector, engineWaitGroup: engineWaitGroup}, nil
}

// Run runs the Gateway
func (g *Gateway) Run() {
	// Wait group
	g.engineWaitGroup.Add(1)
	defer g.engineWaitGroup.Done()

	laddr, err := net.ResolveTCPAddr("tcp", g.bind)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "gateway",
			"bind":   g.bind,
			"error":  err,
			"secure": g.isSecure,
		}).Error("Unable to resolve TCP Address")
		process.SafeExit(1)
		return
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "gateway",
			"bind":   g.bind,
			"error":  err,
			"secure": g.isSecure,
		}).Error("Unable to listen gateway port")
		process.SafeExit(1)
		return
	}

	log.Logger.WithFields(logrus.Fields{
		"prefix": "gateway",
		"bind":   g.bind,
		"secure": g.isSecure,
	}).Info("Started gateway")

	for {

		select {
		case <-g.context.Done():
			listener.Close()
			log.Logger.WithFields(logrus.Fields{
				"prefix": "gateway",
				"secure": g.isSecure,
			}).Info("Stopped server")
			return
		default:
			listener.SetDeadline(time.Now().Add(time.Second))

			conn, err := listener.Accept()
			if err != nil {
				if os.IsTimeout(err) {
					continue
				}

				log.Logger.WithFields(logrus.Fields{
					"prefix": "gateway",
					"error":  err,
					"secure": g.isSecure,
				}).Error("Unable to accept TCP connection")
				continue
			}
			go g.HandleConnection(conn)
		}
	}
}

// Stop stops the gateway server
func (g *Gateway) Stop() {
	g.cancelContextFunc()
}
