package stats

import (
	"context"
	"sync"
	"time"

	"github.com/flexpool/solo/db"
	"github.com/flexpool/solo/log"

	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

const statCollectionPeriodSecs = 60 // Collect stats every minute
const keepStatsForSecs = 86400      // Keep stats for one day

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

// ResetValues deletes (reinitializes) all values from Stats Collector
func (c *Collector) ResetValues() {
	c.Mux.Lock()
	c.PendingStats = make(map[string]PendingStat)
	c.Mux.Unlock()
}

// NewCollector creates a new Stats Collector
func NewCollector(database *db.Database, engineWaitGroup *sync.WaitGroup) *Collector {
	ctx, cancelFunc := context.WithCancel(context.Background())
	c := Collector{
		Context:           ctx,
		ContextCancelFunc: cancelFunc,
		engineWaitGroup:   engineWaitGroup,
	}
	c.ResetValues()
	return &c
}

// Run runs the StatsCollector
func (c *Collector) Run() {
	// Wait group
	c.engineWaitGroup.Add(1)
	defer c.engineWaitGroup.Done()

	// TODO: Collect outdated stats (garbege collection)

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
			if prevCollectionTimestamp != currentCollectionTimestamp {
				time.Sleep(time.Second)
				continue
			}

			prevCollectionTimestamp = currentCollectionTimestamp

			c.Mux.Lock()

			batch := new(leveldb.Batch)
			for workerName, pendingStat := range c.PendingStats {
				timestamp := time.Now().Unix() / statCollectionPeriodSecs * statCollectionPeriodSecs // Get rid of remainder
				effectiveHashrate := float64(pendingStat.ValidShares) * float64(c.ShareDifficulty)
				totalCollectedHashrate += effectiveHashrate
				stat := db.Stat{
					WorkerName:        workerName,
					ValidShareCount:   pendingStat.ValidShares,
					StaleShareCount:   pendingStat.StaleShares,
					InvalidShareCount: pendingStat.InvalidShares,
					ReportedHashrate:  pendingStat.ReportedHashrate,
					EffectiveHashrate: effectiveHashrate,
					IPAddress:         pendingStat.IPAddress,
				}
				db.WriteStatToBatch(batch, stat, timestamp)
			}
			c.Mux.Unlock()

			c.Database.DB.Write(batch, nil)

			log.Logger.WithFields(logrus.Fields{
				"prefix":             "stats",
				"effective-hashrate": totalCollectedHashrate,
			}).Info("Successfully collected data.")
			totalCollectedHashrate = 0

			c.Database.PruneStats(keepStatsForSecs)
		}
	}
}

// Stop function stops stats collector
func (c *Collector) Stop() {
	c.ContextCancelFunc()
}
