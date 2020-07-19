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

package db

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/vmihailenco/msgpack/v5"
)

const statPrefix = "stat__"
const bestSharePrefix = "best__"
const minedBlockPrefix = "block__"

const minedValidSharesKey = "mined_valid_shares"

// Stat represents an interface for a stat DB object
type Stat struct {
	WorkerName        string  `msgpack:"worker_name"`
	ValidShareCount   uint64  `msgpack:"valid_share_count"`
	StaleShareCount   uint64  `msgpack:"stale_share_count"`
	InvalidShareCount uint64  `msgpack:"invalid_share_count"`
	ReportedHashrate  float64 `msgpack:"reported_hashrate"`
	EffectiveHashrate float64 `msgpack:"effective_hashrate"`
	IPAddress         string  `msgpack:"ip_address"`
}

// BestShare represents an interface for a best share DB object
type BestShare struct {
	WorkerName            string  `msgpack:"worker_name"`
	ActualShareDifficulty float64 `msgpack:"actual_share_difficulty"`
}

// Block represents an interface for a block DB object
type Block struct {
	Hash        string  `msgpack:"hash"`
	Number      uint64  `msgpack:"number"`
	Type        string  `msgpack:"type"`
	WorkerName  string  `msgpack:"worker_name"`
	Difficulty  float64 `msgpack:"difficulty"`
	Timestamp   int64   `msgpack:"timestamp"`
	Confirmed   bool    `msgpack:"confirmed"`
	MinedHashes float64 `msgpack:"mined_hashes"`
	RoundTime   int64   `msgpack:"round_time"`
	Luck        float64 `msgpack:"luck"`
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

	key := minedBlockPrefix + block.Hash

	return db.DB.Put([]byte(key), data, nil)
}

// IncrValidShares increments mined valid shares counter (used to precisely calculate luck)
func (db *Database) IncrValidShares() error {
	prevValBytes, _ := db.DB.Get([]byte(minedValidSharesKey), nil)
	prevVal, _ := strconv.ParseUint(string(prevValBytes), 10, 64)
	return db.DB.Put([]byte(minedValidSharesKey), []byte(string(strconv.FormatUint(prevVal+1, 10))), nil)
}

// GetValidSharesThenReset gets the mined valid shares counter and resets it
func (db *Database) GetValidSharesThenReset() (uint64, error) {
	valBytes, err := db.DB.Get([]byte(minedValidSharesKey), nil)
	if err != nil {
		return 0, errors.Wrap(err, "unable to read valid share counter from the db")
	}
	return strconv.ParseUint(string(valBytes), 10, 64)
}

// GetRoundTime returns round time
func (db *Database) GetRoundTime() (int64, error) {
	return 0, errors.New("unimplemented") // Needed to implement querying first
}
