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
	"strconv"
	"sync"
	"time"

	"github.com/flexpool/solo/db"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/nodeapi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/vmihailenco/msgpack/v5"
)

// BlockConfirmationManager is a daemon that confirms blocks by verifying its type (block/uncle/orphan)
type BlockConfirmationManager struct {
	Database              *db.Database
	Node                  *nodeapi.Node
	Context               context.Context
	ContextCancelFunc     context.CancelFunc
	engineWaitGroup       *sync.WaitGroup
	confirmationsRequired uint64
}

// NewBlockConfirmationManager creates a new BlockConfirmationManager instance
func NewBlockConfirmationManager(database *db.Database, engineWaitGroup *sync.WaitGroup, node *nodeapi.Node, confirmationsRequired uint64) *BlockConfirmationManager {
	ctx, contextCancelFunc := context.WithCancel(context.Background())

	return &BlockConfirmationManager{
		Database:              database,
		Context:               ctx,
		ContextCancelFunc:     contextCancelFunc,
		engineWaitGroup:       engineWaitGroup,
		confirmationsRequired: confirmationsRequired,
		Node:                  node,
	}
}

// Run function runs the BlockConfirmationManager
func (b *BlockConfirmationManager) Run() {
	b.engineWaitGroup.Add(1)
	defer b.engineWaitGroup.Done()

	for {
		select {
		case <-b.Context.Done():
			return
		default:

			currentBlockNumber, err := b.Node.BlockNumber()
			if err != nil {
				log.Logger.WithFields(logrus.Fields{
					"prefix": "blockmanager",
					"error":  err.Error(),
				}).Error("Unable to get block number")
			}

			iter := b.Database.DB.NewIterator(util.BytesPrefix([]byte(db.MinedBlockPrefix)), nil)
			for iter.Next() {
				var parsedBlock db.Block
				err := msgpack.Unmarshal(iter.Value(), &parsedBlock)
				if err != nil || (parsedBlock.Type != "block" && parsedBlock.Type != "uncle" && parsedBlock.Type != "orphan") {
					panic(errors.Wrap(err, "Database is corrupted"))
				}
				if !parsedBlock.Confirmed && parsedBlock.Number+b.confirmationsRequired <= currentBlockNumber {
					// Block is unconfirmed, but has enough confirmations.

					blockByNumber, err := b.Node.GetBlockByNumber(parsedBlock.Number)
					if err != nil {
						panic(err)
					}

					if blockByNumber.Hash != parsedBlock.Hash {
						// It may be Uncle or Orphan.
						_, err := b.Node.GetBlockByHash(parsedBlock.Hash)
						if err != nil {
							log.Logger.WithFields(logrus.Fields{
								"prefix":       "blockmanager",
								"lookup-error": err.Error(),
							}).Warn("It seems like " + parsedBlock.Type + " #" + strconv.FormatUint(parsedBlock.Number, 10) + " was orphaned")

							// Error while searching it in blockchain. It is orphaned.
							parsedBlock.Type = "orphan"
						} else {
							parsedBlock.Type = "uncle"
						}
					}

					parsedBlock.Confirmed = true
					blockBytes, _ := msgpack.Marshal(parsedBlock)
					if b.Database.DB.Put(iter.Key(), blockBytes, nil) != nil {
						log.Logger.WithFields(logrus.Fields{
							"prefix": "blockmanager",
							"error":  err,
						}).Error("Unable to write confirmed block")
					}

					log.Logger.WithFields(logrus.Fields{
						"prefix": "blockmanager",
						"type":   parsedBlock.Type,
						"number": parsedBlock.Number,
						"hash":   parsedBlock.Hash,
					}).Info("✅ Block confirmed")
				}
			}
			iter.Release()
			time.Sleep(1 * time.Second)
		}
	}
}

// Stop function stops the BlockConfirmationManager
func (b *BlockConfirmationManager) Stop() {
	b.ContextCancelFunc()
}
