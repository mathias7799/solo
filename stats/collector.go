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

package stats

import (
	"context"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/flexpool/solo/db"
	"github.com/flexpool/solo/log"

	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

const statCollectionPeriodSecs = 600 // Collect stats every 10 minutes
const keepStatsForSecs = 86400       // Keep stats for one day

// Collector is a stat collection daemon struct
type Collector struct {
	// map[<worker-name>]PendingStat
	PendingStats map[string]PendingStat

	ShareDifficulty   uint64
	Database          *db.Database
	Context           context.Context
	ContextCancelFunc context.CancelFunc
	Mux               sync.Mutex
	engineWaitGroup   *sync.WaitGroup
}

// Init initializes the Collector
func (c *Collector) Init() {
	c.Mux.Lock()
	c.PendingStats = make(map[string]PendingStat)
	c.Mux.Unlock()
}

// Clear deletes all keys (and values) from the c.PendingStats
func (c *Collector) Clear() {
	c.Mux.Lock()
	for k := range c.PendingStats {
		delete(c.PendingStats, k)
	}
	c.Mux.Unlock()
}

// NewCollector creates a new Stats Collector
func NewCollector(database *db.Database, engineWaitGroup *sync.WaitGroup, shareDifficulty uint64) *Collector {
	ctx, cancelFunc := context.WithCancel(context.Background())
	c := Collector{
		Context:           ctx,
		ContextCancelFunc: cancelFunc,
		engineWaitGroup:   engineWaitGroup,
		Database:          database,
		ShareDifficulty:   shareDifficulty,
	}
	c.Init()
	return &c
}

// Run runs the StatsCollector
func (c *Collector) Run() {
	// Wait group
	c.engineWaitGroup.Add(1)
	defer c.engineWaitGroup.Done()

	prevCollectionTimestamp := time.Now().Unix() / statCollectionPeriodSecs * statCollectionPeriodSecs

	log.Logger.WithFields(logrus.Fields{
		"prefix": "stats",
	}).Info("Started Stats Collector")

	var totalCollectedHashrate float64

	for {
		select {
		case <-c.Context.Done():
			log.Logger.WithFields(logrus.Fields{
				"prefix": "stats",
			}).Info("Stopped Stats Collector")
			return
		default:
			currentCollectionTimestamp := time.Now().Unix() / statCollectionPeriodSecs * statCollectionPeriodSecs // Get rid of remainder
			if prevCollectionTimestamp == currentCollectionTimestamp {
				time.Sleep(time.Second)
				continue
			}

			prevCollectionTimestamp = currentCollectionTimestamp

			c.Mux.Lock()

			batch := new(leveldb.Batch)
			pendingTotalStat := db.TotalStat{}
			timestamp := time.Now().Unix() / statCollectionPeriodSecs * statCollectionPeriodSecs // Get rid of remainder

			for workerName, pendingStat := range c.PendingStats {
				effectiveHashrate := float64(pendingStat.ValidShares) * float64(c.ShareDifficulty)
				totalCollectedHashrate += effectiveHashrate / statCollectionPeriodSecs
				stat := db.Stat{
					WorkerName:        workerName,
					ValidShareCount:   pendingStat.ValidShares,
					StaleShareCount:   pendingStat.StaleShares,
					InvalidShareCount: pendingStat.InvalidShares,
					ReportedHashrate:  pendingStat.ReportedHashrate,
					EffectiveHashrate: effectiveHashrate,
					IPAddress:         pendingStat.IPAddress,
				}

				pendingTotalStat.ValidShareCount += pendingStat.ValidShares
				pendingTotalStat.StaleShareCount += pendingStat.StaleShares
				pendingTotalStat.InvalidShareCount += pendingStat.InvalidShares
				pendingTotalStat.EffectiveHashrate += effectiveHashrate
				pendingTotalStat.ReportedHashrate += pendingStat.ReportedHashrate
				pendingTotalStat.WorkerCount++

				db.WriteStatToBatch(batch, stat, timestamp)
			}

			c.Mux.Unlock()
			c.Clear()

			db.WriteTotalStatToBatch(batch, pendingTotalStat, timestamp)

			c.Database.DB.Write(batch, nil)

			log.Logger.WithFields(logrus.Fields{
				"prefix":             "stats",
				"effective-hashrate": humanize.SIWithDigits(totalCollectedHashrate, 2, "H/s"),
			}).Info("Collected data")
			totalCollectedHashrate = 0

			c.Database.GetAndWriteCachedValues()
			c.Database.PruneStats(keepStatsForSecs)
		}
	}
}

// Stop function stops stats collector
func (c *Collector) Stop() {
	c.ContextCancelFunc()
}
