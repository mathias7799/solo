package db

import (
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/vmihailenco/msgpack/v5"
)

const statPrefix = "stat__"
const bestSharePrefix = "best__"
const minedBlockPrefix = "block__"

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

// Block represents an interface for a block DB object
type Block struct {
	Hash                  string
	Number                uint64
	Type                  string
	Worker                string
	Difficulty            string
	Timestamp             string
	Confirmed             bool
	MinedHashes           float64
	RoundTime             int64
	Luck                  uint64
	BlockReward           *big.Int
	BlockFees             *big.Int
	UncleInclusionRewards *big.Int
	TotalRewards          *big.Int
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

// PruneStats removes data older than
func (db *Database) PruneStats(deleteDataOlderThanSecs int64) {
	iter := db.DB.NewIterator(util.BytesPrefix([]byte(statPrefix)), nil)

	deleteWithTimestampLowerThan := time.Now().Unix() - deleteDataOlderThanSecs

	for iter.Next() {
		key := iter.Key()
		keySplitted := strings.Split(string(key), "_")
		timestamp, err := strconv.ParseInt(keySplitted[len(keySplitted)-1], 10, 64)
		if err != nil {
			panic("Database corruption")
		}

		if timestamp < deleteWithTimestampLowerThan {
			db.DB.Delete(key, nil)
		}
	}
}

// WriteMinedBlock writes mined block to the database
func (db *Database) WriteMinedBlock(block Block) error {
	data, err := msgpack.Marshal(block)
	if err != nil {
		fmt.Println("DEBUG Print on WriteMinedBlock (don't forget to remove panic)")
		panic(err)
	}

	key := bestSharePrefix + minedBlockPrefix + block.Hash

	return db.DB.Put([]byte(key), data, nil)
}
