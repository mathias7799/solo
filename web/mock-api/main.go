package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const apiPrefix = "/api/v1"

func getCurrent1MinTimestamp() int64 {
	return time.Now().Unix() / 600 * 600 // Get rid of the remainder
}

type history struct {
	Timestamp int64   `json:"timestamp"`
	Effective float64 `json:"effectiveHashrate"`
	Reported  float64 `json:"reportedHashrate"`
	Valid     float64 `json:"validShares"`
	Stale     float64 `json:"staleShares"`
	Invalid   float64 `json:"invalidShares"`
}

const shareDifficulty float64 = 4000000000

func genTotalHistory() []history {
	var totalHistory []history

	current1MinTimestamp := getCurrent1MinTimestamp()

	for i := int64(144); i != 0; i-- {
		validShares := float64(200 + rand.Intn(50))

		staleShares := float64(rand.Intn(20))
		var invalidShares float64
		if rand.Intn(100) < 30 {
			invalidShares = float64(rand.Intn(2))
		}

		effectiveHashrate := validShares * shareDifficulty / 600

		totalHistory = append(totalHistory, history{
			Timestamp: current1MinTimestamp - i*600,
			Effective: effectiveHashrate,
			Reported:  effectiveHashrate * float64(0.9),
			Valid:     validShares,
			Stale:     staleShares,
			Invalid:   invalidShares,
		})
	}

	return totalHistory
}

func getSI(number float64) (float64, string) {
	if number < 1000 {
		return 1, ""
	}
	symbols := "kMGTPEZY"
	symbolsLen := len(symbols)
	i := 1
	for {
		number /= 1000
		if number < 1000 || i == symbolsLen-1 {
			return math.Pow(1000, float64(i)), string(symbols[i-1])
		}
		i++
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))

	totalHistory := genTotalHistory()

	effectiveHashrate := totalHistory[144-1].Effective
	reportedHashrate := totalHistory[144-1].Reported

	validShares := totalHistory[15].Valid
	staleShares := totalHistory[15].Stale
	invalidShares := totalHistory[15].Invalid

	efficiency := validShares / (validShares + staleShares + invalidShares) * 100

	workersOnline := rand.Intn(5)
	workersOffline := rand.Intn(5)
	fmt.Println(workersOnline)

	balance := rand.Uint64()
	chainID := 1

	var averageHashrate float64
	for _, d := range totalHistory {
		averageHashrate += d.Effective / 144
		validShares += d.Valid
		staleShares += d.Stale
		invalidShares += d.Invalid
	}

	siDiv, siChar := getSI(averageHashrate)

	r.GET(apiPrefix+"/currentBlock", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": gin.H{
				"blockNumber": 12345678,
				"syncing": gin.H{
					"status":       false,
					"currentBlock": 12345678,
					"targetBlock":  12345678,
				},
			},
			"error": nil,
		})
	})

	r.GET(apiPrefix+"/coinbaseBalance", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": balance,
			"error":  nil,
		})
	})

	r.GET(apiPrefix+"/history", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": totalHistory,
			"error":  nil,
		})
	})

	r.GET(apiPrefix+"/headerStats", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": gin.H{
				"workersOnline":   workersOnline,
				"workersOffline":  workersOffline,
				"coinbaseBalance": balance,
				"efficiency":      efficiency,
				"chainId":         chainID,
			},
			"error": nil,
		})
	})

	r.GET(apiPrefix+"/stats", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": gin.H{
				"hashrate": gin.H{
					"effective": effectiveHashrate,
					"reported":  reportedHashrate,
					"average":   averageHashrate,
				},
				"shares": gin.H{
					"valid":   validShares,
					"stale":   staleShares,
					"invalid": invalidShares,
				},
				"si": gin.H{
					"div":  siDiv,
					"char": siChar,
				},
			},
			"error": nil,
		})
	})

	r.Run("127.0.0.1:8000")
}
