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
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/flexpool/solo/log"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/vmihailenco/msgpack/v5"
)

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

// TotalStat represents an interface for a summarized stat DB object
type TotalStat struct {
	ValidShareCount   uint64 `msgpack:"valid_share_count"`
	StaleShareCount   uint64 `msgpack:"stale_share_count"`
	InvalidShareCount uint64 `msgpack:"invalid_share_count"`
	WorkerCount       uint64 `msgpack:"worker_count"`
}

// BestShare represents an interface for a best share DB object
type BestShare struct {
	WorkerName            string  `msgpack:"worker_name"`
	ActualShareDifficulty float64 `msgpack:"actual_share_difficulty"`
	Timestamp             int64   `msgpack:"timestamp"`
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
	key := StatPrefix + stat.WorkerName + "_" + strconv.FormatInt(timestamp, 10)
	batch.Put([]byte(key), data)
}

// WriteTotalStatToBatch writes worker stat object to the LevelDB batch
func WriteTotalStatToBatch(batch *leveldb.Batch, stat TotalStat, timestamp int64) {
	data, _ := msgpack.Marshal(stat)
	key := TotalStatPrefix + "_" + strconv.FormatInt(timestamp, 10)
	batch.Put([]byte(key), data)
}

// PruneStats removes data older than
func (db *Database) PruneStats(deleteDataOlderThanSecs int64) {
	iter := db.DB.NewIterator(util.BytesPrefix([]byte(StatPrefix)), nil)

	deleteWithTimestampLowerThan := time.Now().Unix() - deleteDataOlderThanSecs

	for iter.Next() {
		key := iter.Key()
		keySplitted := strings.Split(string(key), "_")
		timestamp, err := strconv.ParseInt(keySplitted[len(keySplitted)-1], 10, 64)
		if err != nil {
			panic(errors.Wrap(err, "Database is corrupted"))
		}

		if timestamp < deleteWithTimestampLowerThan {
			db.DB.Delete(key, nil)
		}
	}

	iter.Release()
}

// WriteMinedBlock writes mined block to the database
func (db *Database) WriteMinedBlock(block Block) error {
	data, _ := msgpack.Marshal(block)
	key := MinedBlockPrefix + block.Hash
	return db.DB.Put([]byte(key), data, nil)
}

// WriteBestShare writes best share  to the database
func (db *Database) WriteBestShare(bestShare BestShare, timestamp int64) error {
	data, _ := msgpack.Marshal(bestShare)
	key := BestSharePrefix + bestShare.WorkerName + "_" + strconv.FormatInt(timestamp, 10) + "_" + strconv.FormatUint(rand.Uint64(), 16)
	return db.DB.Put([]byte(key), data, nil)
}

// IncrValidShares increments mined valid shares counter (used to precisely calculate luck)
func (db *Database) IncrValidShares() error {
	prevValBytes, _ := db.DB.Get([]byte(MinedValidSharesKey), nil)
	prevVal, _ := strconv.ParseUint(string(prevValBytes), 10, 64)
	return db.DB.Put([]byte(MinedValidSharesKey), []byte(string(strconv.FormatUint(prevVal+1, 10))), nil)
}

// GetValidSharesThenReset gets the mined valid shares counter and resets it
func (db *Database) GetValidSharesThenReset() (uint64, error) {
	valBytes, err := db.DB.Get([]byte(MinedValidSharesKey), nil)
	if err != nil {
		return 0, errors.Wrap(err, "unable to read valid share counter from the db")
	}
	db.DB.Delete([]byte(MinedValidSharesKey), nil)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "db",
			"error":  err,
		}).Error("Unable to delete mined valid shares counter")
	}
	return strconv.ParseUint(string(valBytes), 10, 64)
}

// GetRoundTime returns round time
func (db *Database) GetRoundTime() int64 {
	blocks := db.GetBlocksUnsorted()

	if len(blocks) > 0 {
		var latestBlockTimestamp int64 = time.Now().Unix()
		for _, block := range blocks {
			if latestBlockTimestamp > block.Timestamp {
				latestBlockTimestamp = block.Timestamp
			}
		}
		return time.Now().Unix() - latestBlockTimestamp
	}

	// Search for first best share
	bestShares := db.GetUnsortedBestShares()

	var earliestBestShareTimestamp int64 = time.Now().Unix()

	for _, bestShare := range bestShares {
		if earliestBestShareTimestamp > bestShare.Timestamp {
			earliestBestShareTimestamp = bestShare.Timestamp
		}
	}

	return time.Now().Unix() - earliestBestShareTimestamp
}

// GetBlocksUnsorted returns the unsorted blocks from the Database
func (db *Database) GetBlocksUnsorted() []Block {
	var blocks []Block
	iter := db.DB.NewIterator(util.BytesPrefix([]byte(MinedBlockPrefix)), nil)
	for iter.Next() {
		blockHash := strings.Replace(string(iter.Key()), MinedBlockPrefix, "", 1)
		if hashBytes, err := hex.DecodeString(blockHash[2:]); err != nil || len(hashBytes) > 32 {
			panic("Database is corrupted")
		}

		var parsedBlock Block

		if err := msgpack.Unmarshal(iter.Value(), &parsedBlock); err != nil {
			panic(errors.Wrap(err, "Database is corrupted"))
		}

		if parsedBlock.Hash != blockHash {
			panic("Database is corruped (Blockhash Key: \"" + blockHash + "\", Actual: \"" + parsedBlock.Hash + "\")")
		}

		blocks = append(blocks, parsedBlock)
	}

	iter.Release()
	return blocks
}

// GetUnsortedBestShares returns the unsorted best shares from the Database
func (db *Database) GetUnsortedBestShares() []BestShare {
	var bestShares []BestShare
	iter := db.DB.NewIterator(util.BytesPrefix([]byte(BestSharePrefix)), nil)
	for iter.Next() {
		timestampString := strings.Split(strings.Split(strings.Replace(string(iter.Key()), BestSharePrefix, "", 1), "__")[0], "_")[1]
		timestamp, err := strconv.ParseInt(timestampString, 10, 64)
		if err != nil {
			panic(errors.Wrap(err, "Database is corrupted"))
		}

		var parsedBestShare BestShare

		if err := msgpack.Unmarshal(iter.Value(), &parsedBestShare); err != nil {
			panic(errors.Wrap(err, "Database is corrupted"))
		}

		if timestamp != parsedBestShare.Timestamp {
			panic(fmt.Sprintf("Database is corrupted (Timestamp key: %v, Actual: %v)", timestamp, parsedBestShare.Timestamp))
		}

		bestShares = append(bestShares, parsedBestShare)

	}
	iter.Release()
	return bestShares
}
