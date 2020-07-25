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

package web

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"sync"

	"github.com/flexpool/solo/db"
	"github.com/flexpool/solo/gateway"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/nodeapi"
	"github.com/flexpool/solo/process"
	"github.com/flexpool/solo/utils"

	"github.com/sirupsen/logrus"
)

// h is a Simple shortcut to map[string]interface{}
type h map[string]interface{}

// Server is a RESTful API & Front End App server instance
type Server struct {
	httpServer      *http.Server
	database        *db.Database
	workmanager     *gateway.WorkManager
	node            *nodeapi.Node
	engineWaitGroup *sync.WaitGroup
	shuttingDown    bool
}

// APIResponse is an interface to APIResponse
type APIResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

// MarshalAPIResponse function marshals APIResponse struct
func MarshalAPIResponse(resp APIResponse) []byte {
	data, _ := json.Marshal(resp)
	return data
}

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

// NewServer creates new Server instance
func NewServer(db *db.Database, node *nodeapi.Node, engineWaitGroup *sync.WaitGroup, workmanager *gateway.WorkManager, bind string) *Server {
	mux := http.NewServeMux()

	server := Server{
		database:        db,
		node:            node,
		workmanager:     workmanager,
		engineWaitGroup: engineWaitGroup,
	}

	mux.HandleFunc("/api/v1/currentBlock", func(w http.ResponseWriter, r *http.Request) {

		// DEBUG ONLY !!!
		w.Header().Set("Access-Control-Allow-Origin", "*")

		currentBlock, err := server.node.BlockNumber()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(MarshalAPIResponse(APIResponse{
			Result: currentBlock,
			Error:  err,
		}))
	})

	mux.HandleFunc("/api/v1/stats", func(w http.ResponseWriter, r *http.Request) {
		// DEBUG ONLY !!!
		w.Header().Set("Access-Control-Allow-Origin", "*")

		currentTotalStats, err := server.database.GetTotalStatsByTimestamp(utils.GetCurrent10MinTimestamp())
		totalShares, err := server.database.GetTotalShares()

		averageEffective := db.GetTotalAverageHashrate()
		averageEffectiveFloat, _ := big.NewFloat(0).SetInt(averageEffective).Float64()

		siDiv, siChar := utils.GetSI(averageEffectiveFloat)

		data, _ := json.Marshal(h{
			"result": h{
				"hashrate": h{
					"effective": currentTotalStats.EffectiveHashrate,
					"reported":  currentTotalStats.ReportedHashrate,
					"average":   averageEffective,
				},
				"shares": h{
					"valid":   totalShares.ValidShares,
					"stale":   totalShares.StaleShares,
					"invalid": totalShares.InvalidShares,
				},
				"si": h{
					"div":  siDiv,
					"char": siChar,
				},
			},
			"error": err,
		})
		w.Write(data)
	})

	server.httpServer = &http.Server{
		Addr:    bind,
		Handler: mux,
	}

	return &server
}

// Run function runs the Server
func (a *Server) Run() {
	a.engineWaitGroup.Add(1)

	err := a.httpServer.ListenAndServe()

	if !a.shuttingDown {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "web",
			"error":  err.Error(),
		}).Error("API Server shut down unexpectedly")
		a.engineWaitGroup.Done()
		process.SafeExit(1)
	}

	a.engineWaitGroup.Done()
}

// Stop function stops the Server
func (a *Server) Stop() {
	a.shuttingDown = true
	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}
