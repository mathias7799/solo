package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
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

type worker struct {
	Effective float64 `json:"effectiveHashrate"`
	Reported  float64 `json:"reportedHashrate"`
	Valid     float64 `json:"validShares"`
	Stale     float64 `json:"staleShares"`
	Invalid   float64 `json:"invalidShares"`
	LastSeen  int64   `json:"lastSeen"`
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
			Reported:  effectiveHashrate * float64(9+rand.Intn(6)) / 12,
			Valid:     validShares,
			Stale:     staleShares,
			Invalid:   invalidShares,
		})
	}

	return totalHistory
}

func genRandomWorker(online bool) worker {
	validShares := float64(50 + rand.Intn(10))
	staleShares := float64(rand.Intn(5))
	invalidShares := float64(rand.Intn(2))

	effectiveHashrate := validShares * shareDifficulty / 600
	lastSeen := int64(rand.Intn(600))
	if !online {
		lastSeen += 600
	}

	return worker{
		Effective: effectiveHashrate,
		Reported:  effectiveHashrate * float64(9+rand.Intn(6)) / 12,
		Valid:     validShares,
		Stale:     staleShares,
		Invalid:   invalidShares,
		LastSeen:  lastSeen,
	}
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

	var workers = make(map[string]worker)
	for i := 0; i < workersOffline; i++ {
		workers[strconv.Itoa(rand.Intn(10))+"-off"] = genRandomWorker(false)
	}
	for i := 0; i < workersOnline; i++ {
		workers[strconv.Itoa(rand.Intn(10))] = genRandomWorker(true)
	}

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
		rand.Seed(time.Now().UnixNano())
		fmt.Println("param", c.Param("workerName"), c.Params)
		fmt.Println(c.Request.URL.Query())
		if len(c.Request.URL.Query()["workerName"]) == 0 || c.Request.URL.Query()["workerName"][0] == "" {
			c.JSON(200, gin.H{
				"result": totalHistory,
				"error":  nil,
			})
		} else {
			c.JSON(200, gin.H{
				"result": genTotalHistory(),
				"error":  nil,
			})
		}
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
		var localEffective float64
		var localReported float64
		var localAverage float64

		var localValid uint64
		var localStale uint64
		var localInvalid uint64

		var localSiChar string
		var localSiDiv float64 = 1

		if len(c.Request.URL.Query()["workerName"]) == 0 || c.Request.URL.Query()["workerName"][0] == "" {
			localEffective = effectiveHashrate
			localReported = reportedHashrate
			localAverage = averageHashrate

			localValid = uint64(validShares)
			localStale = uint64(staleShares)
			localInvalid = uint64(invalidShares)

			localSiDiv, localSiChar = siDiv, siChar
		} else {
			workerName := c.Request.URL.Query()["workerName"][0]
			worker := workers[workerName]

			localEffective = worker.Effective
			localReported = worker.Reported
			localAverage = worker.Effective // There's no average property in worker struct

			localValid = uint64(worker.Valid)
			localStale = uint64(worker.Stale)
			localInvalid = uint64(worker.Invalid)

			localSiDiv, localSiChar = getSI(localEffective)
		}
		c.JSON(200, gin.H{
			"result": gin.H{
				"hashrate": gin.H{
					"effective": localEffective,
					"reported":  localReported,
					"average":   localAverage,
				},
				"shares": gin.H{
					"valid":   localValid,
					"stale":   localStale,
					"invalid": localInvalid,
				},
				"si": gin.H{
					"div":  localSiDiv,
					"char": localSiChar,
				},
			},
			"error": nil,
		})
	})

	r.GET(apiPrefix+"/workers", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": workers,
			"error":  nil,
		})
	})

	r.Run("127.0.0.1:8000")
}
