package db

import (
	"math/rand"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/vmihailenco/msgpack/v5"
)

const statPrefix = "stat__"
const bestSharePrefix = "best__"

// Stat represents an interface for a stat DB object
type Stat struct {
	WorkerName        string
	ValidShareCount   uint64
	StaleShareCount   uint64
	InvalidShareCount uint64
	ReportedHashrate  float64
	EffectiveHashrate float64
	IPAddress         string
}

// BestShare represents an interface for a best share DB object
type BestShare struct {
	WorkerName            string
	ActualShareDifficulty float64
}

// WriteStatToBatch writes worker stat object to the LevelDB batch
func WriteStatToBatch(batch *leveldb.Batch, stat Stat, timestamp int64) {
	data, _ := msgpack.Marshal(stat)
	key := statPrefix + stat.WorkerName + "_" + strconv.FormatInt(timestamp, 10)
	batch.Put([]byte(key), data)
}

// WriteBestShareToBatch writes best share object to the LevelDB batch
func WriteBestShareToBatch(batch *leveldb.Batch, bestShare BestShare, timestamp int64) {
	data, _ := msgpack.Marshal(bestShare)
	key := bestSharePrefix + bestShare.WorkerName + "_" + strconv.FormatInt(timestamp, 10) + "_" + strconv.FormatUint(rand.Uint64(), 16)
	batch.Put([]byte(key), data)
}
