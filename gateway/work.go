package gateway

import (
	"errors"
	"math/big"

	"github.com/dustin/go-humanize"
	"github.com/ethereum/go-ethereum/common"
	"github.com/flexpool/ethash-go"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/utils"
	"github.com/sirupsen/logrus"
)

const (
	shareValid   = 0
	shareStale   = 1
	shareInvalid = 2
)

// Block is ethash validation block structure
type Block struct {
	difficulty  *big.Int
	hashNoNonce common.Hash
	nonce       uint64
	mixDigest   common.Hash
	number      uint64
}

// Difficulty Returns block's difficulty
func (b Block) Difficulty() *big.Int { return b.difficulty }

// HashNoNonce Returns block's hash
func (b Block) HashNoNonce() common.Hash { return b.hashNoNonce }

// Nonce Returns block's nonce
func (b Block) Nonce() uint64 { return b.nonce }

// MixDigest Returns block's mix digest
func (b Block) MixDigest() common.Hash { return b.mixDigest }

// NumberU64 Returns block's number
func (b Block) NumberU64() uint64 { return b.number }

var hasher = ethash.New()

func (g *Gateway) validateShare(workerWork []string, workerName string) (int8, error) {
	// workerName is required to know who mined the block, if there share mines it

	g.parentWorkReceiver.workHistory.Mux.Lock()
	fullWork, ok := g.parentWorkReceiver.workHistory.Map[workerWork[0]]
	g.parentWorkReceiver.workHistory.Mux.Unlock()
	if !ok {
		// Work was not requested, or is older than 8 blocks
		return shareInvalid, errors.New("Work is outdated, or not requested")
	}

	var isStale bool

	blockNumber := utils.MustSoftHexToUint64(fullWork[3])
	if fullWork[3] != g.parentWorkReceiver.lastWork[3] {
		isStale = true
	}

	share := Block{
		difficulty:  g.parentWorkReceiver.shareTargetBigInt,
		hashNoNonce: common.HexToHash(fullWork[0]),
		nonce:       utils.MustSoftHexToUint64(workerWork[0]),
		mixDigest:   common.HexToHash(fullWork[2]),
		number:      blockNumber,
	}

	shareIsValid, actualTarget := hasher.Verify(share)

	if shareIsValid {
		block := share
		block.difficulty = utils.HexStrToBigInt(fullWork[2])

		blockValid, _ := hasher.Verify(block)

		if blockValid {
			// TODO: Block mined, submit it.
		}

		if g.parentWorkReceiver.BestShareTarget.Cmp(actualTarget) == 1 {
			float64ActualDifficulty, _ := big.NewFloat(0).SetInt(big.NewInt(0).Div(utils.BigMax256bit, actualTarget)).Float64()
			log.Logger.WithFields(logrus.Fields{
				"prefix":            "gateway",
				"actual-difficulty": humanize.SIWithDigits(float64ActualDifficulty, 2, "H"),
			}).Info("New best share")
			if !blockValid {
				g.parentWorkReceiver.BestShareTarget = actualTarget
			} else {
				g.parentWorkReceiver.BestShareTarget = utils.BigMax256bit
			}
		}

		if isStale {
			return shareStale, nil
		}
		return shareValid, nil
	}

	return shareInvalid, nil
}
