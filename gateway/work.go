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
	"errors"
	"math/big"
	"time"

	"github.com/flexpool/ethash-go"
	"github.com/flexpool/solo/db"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/types"
	"github.com/flexpool/solo/utils"

	"github.com/dustin/go-humanize"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

// Block is ethash validation block structure
type Block struct {
	target      *big.Int
	hashNoNonce common.Hash
	nonce       uint64
	mixDigest   common.Hash
	number      uint64
}

// TargetDifficulty Returns block's difficulty
func (b Block) TargetDifficulty() *big.Int { return b.target }

// HashNoNonce Returns block's hash
func (b Block) HashNoNonce() common.Hash { return b.hashNoNonce }

// Nonce Returns block's nonce
func (b Block) Nonce() uint64 { return b.nonce }

// MixDigest Returns block's mix digest
func (b Block) MixDigest() common.Hash { return b.mixDigest }

// NumberU64 Returns block's number
func (b Block) NumberU64() uint64 { return b.number }

var hasher = ethash.New()

func (g *Gateway) submitBlock(submittedWork []string, blockNumber uint64, workerName string, actualTarget *big.Int) {
	// Submitting the work first
	status, err := g.parentWorkManager.Node.SubmitWork(submittedWork)

	acutalShareDifficulty, _ := big.NewFloat(0).SetInt(big.NewInt(0).Div(utils.BigMax256bit, actualTarget)).Float64()
	log.Logger.WithFields(logrus.Fields{
		"prefix":       "gateway",
		"nonce":        submittedWork[0],
		"block-number": blockNumber,
		"worker":       workerName,
		"actual-diff":  humanize.SIWithDigits(acutalShareDifficulty, 2, "H"),
	}).Info("⛏ Mined potential block")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Unable to submit mined block!")
		return
	}

	if !status {
		log.Logger.WithFields(logrus.Fields{
			"prefix":       "gateway",
			"block-number": blockNumber,
		}).Error("Submitted block marked as invalid!")
		return
	}

	harvestedBlock, uncleParent, err := g.parentWorkManager.Node.HarvestBlockByNonce(submittedWork[0], blockNumber)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"prefix":       "gateway",
			"block-number": blockNumber,
			"error":        err,
		}).Error("Mined block wasn't found in blockchain!")
		return
	}

	blockType := "block"
	if uncleParent > 0 {
		blockType = "uncle"
	}

	log.Logger.WithFields(logrus.Fields{
		"prefix":       "gateway",
		"block-number": blockNumber,
		"type":         blockType,
		"hash":         harvestedBlock.Hash,
	}).Info("⚡️ Submitted block found in blockchain")

	sharesMined, err := g.statsCollector.Database.GetValidSharesThenReset()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "gateway",
			"error":  err.Error(),
		}).Warn("Unable to get valid share counter")
	}

	roundTime, err := g.statsCollector.Database.GetRoundTime()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "gateway",
			"error":  err.Error(),
		}).Error("Unable to round time")
	}

	difficulty := float64(utils.MustSoftHexToUint64(harvestedBlock.Difficulty))
	hashesMined := float64(sharesMined) * float64(g.parentWorkManager.shareDiff)

	// Writing found block to DB
	g.statsCollector.Database.WriteMinedBlock(db.Block{
		Hash:        harvestedBlock.Hash,
		Number:      utils.MustSoftHexToUint64(harvestedBlock.Number),
		Type:        blockType,
		WorkerName:  workerName,
		Difficulty:  difficulty,
		Timestamp:   int64(utils.MustSoftHexToUint64(harvestedBlock.Timestamp)),
		Confirmed:   false,
		MinedHashes: hashesMined,
		RoundTime:   roundTime,
		Luck:        difficulty / hashesMined,
	})
}

func (g *Gateway) validateShare(submittedWork []string, workerName string) (types.ShareType, error) {
	// workerName is required to know who mined the block, if there share mines it

	g.parentWorkManager.workHistory.Mux.Lock()
	fullWork, ok := g.parentWorkManager.workHistory.Map[submittedWork[1]]
	g.parentWorkManager.workHistory.Mux.Unlock()
	if !ok {
		// Work was not requested, or is older than 8 blocks
		return types.ShareInvalid, errors.New("Work is outdated, or not requested")
	}

	var isStale bool

	blockNumber := utils.MustSoftHexToUint64(fullWork[3])
	if fullWork[3] != g.parentWorkManager.lastWork[3] {
		isStale = true
	}

	share := Block{
		target:      g.parentWorkManager.shareTargetBigInt,
		hashNoNonce: common.HexToHash(fullWork[0]),
		nonce:       utils.MustSoftHexToUint64(submittedWork[0]),
		mixDigest:   common.HexToHash(submittedWork[2]),
		number:      blockNumber,
	}

	shareIsValid, actualTarget := hasher.Verify(share)

	if shareIsValid {
		if utils.HexStrToBigInt(fullWork[2]).Cmp(actualTarget) > 0 {
			go g.submitBlock(submittedWork, blockNumber, workerName, actualTarget)
			g.parentWorkManager.BestShareTarget = utils.BigMax256bit
		} else {
			if g.parentWorkManager.BestShareTarget.Cmp(actualTarget) == 1 {
				float64ActualDifficulty, _ := big.NewFloat(0).SetInt(big.NewInt(0).Div(utils.BigMax256bit, actualTarget)).Float64()
				log.Logger.WithFields(logrus.Fields{
					"prefix":      "gateway",
					"actual-diff": humanize.SIWithDigits(float64ActualDifficulty, 2, "H"),
				}).Info("New best share")
				err := g.statsCollector.Database.WriteBestShare(db.BestShare{
					WorkerName:            workerName,
					ActualShareDifficulty: float64ActualDifficulty,
					Timestamp:             time.Now().Unix(),
				}, time.Now().Unix())
				if err != nil {
					log.Logger.WithFields(logrus.Fields{
						"prefix": "gateway",
					}).Error("Unable to write best share to DB")
				}
				g.parentWorkManager.BestShareTarget = actualTarget
			}
		}

		if isStale {
			return types.ShareStale, nil
		}
		return types.ShareValid, nil
	}

	return types.ShareInvalid, nil
}

func (g *Gateway) submitShare(work []string, workerName string) (types.ShareType, error) {
	return g.validateShare(work, workerName)
}
