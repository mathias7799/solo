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
	"math/big"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
)

// TotalShares represents an interface to total shares db object
type TotalShares struct {
	ValidShares   uint64 `msgpack:"valid_shares"`
	StaleShares   uint64 `msgpack:"stale_shares"`
	InvalidShares uint64 `msgpack:"invalidShares"`
}

// GetAndWriteCachedValues gets (calculates) and writes cached values to the db
func (db *Database) GetAndWriteCachedValues() error {
	history, err := db.GetTotalHistory()
	if err != nil {
		return err
	}

	var effectiveNoZeroes []float64
	var totalValidShares uint64
	var totalStaleShares uint64
	var totalInvalidShares uint64

	for _, item := range history {
		if item.EffectiveHashrate != 0 {
			effectiveNoZeroes = append(effectiveNoZeroes, item.EffectiveHashrate)
		}
		totalValidShares += item.ValidShareCount
		totalStaleShares += item.StaleShareCount
		totalInvalidShares += item.InvalidShareCount
	}

	effectiveNoZeroesLen := float64(len(effectiveNoZeroes))

	avgHashrate := big.NewFloat(0)

	for _, item := range effectiveNoZeroes {
		avgHashrate.Set(big.NewFloat(0).Add(avgHashrate, big.NewFloat(item/effectiveNoZeroesLen))) // 144 is the length of history
	}

	avgHashrateInt, _ := avgHashrate.Int(nil)

	err = db.DB.Put([]byte(AverageTotalHashrateKey), avgHashrateInt.Bytes(), nil)
	if err != nil {
		return err
	}
	data, _ := msgpack.Marshal(TotalShares{
		ValidShares:   totalValidShares,
		StaleShares:   totalStaleShares,
		InvalidShares: totalInvalidShares,
	})

	return db.DB.Put([]byte(TotalSharesKey), data, nil)
}

// GetTotalAverageHashrate returns total average hashrate which was put by GetAndWriteCachedValues
func (db *Database) GetTotalAverageHashrate() *big.Int {
	v, _ := db.DB.Get([]byte(AverageTotalHashrateKey), nil)
	return big.NewInt(0).SetBytes(v)
}

// GetTotalShares returns total amount of shares which was put by GetAndWriteCachedValues
func (db *Database) GetTotalShares() (TotalShares, error) {
	data, err := db.DB.Get([]byte(TotalSharesKey), nil)
	if err != nil {
		return TotalShares{}, err
	}
	var totalShares TotalShares
	err = msgpack.Unmarshal(data, &totalShares)
	if err != nil {
		panic(errors.Wrap(err, "Database is corrupted"))
	}
	return totalShares, nil
}
