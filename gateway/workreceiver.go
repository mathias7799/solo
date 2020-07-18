package gateway

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"sync"

	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/utils"

	"github.com/sirupsen/logrus"
)

// OrderedWorkMap is used to store work history, and have an ability to prune unneeded work
type OrderedWorkMap struct {
	Map   map[string][]string
	Order []string
	Mux   sync.Mutex
}

// Init initializes the OrderedWorkMap
func (o *OrderedWorkMap) Init() {
	o.Map = make(map[string][]string)
}

// Append appends new work to the OrderedWorkMap
func (o *OrderedWorkMap) Append(headerHash string, work []string) {
	o.Mux.Lock()
	o.Map[headerHash] = work
	o.Order = append(o.Order, headerHash)
	o.Mux.Unlock()
}

// Shift removes the first OrderedWorkMap key
func (o *OrderedWorkMap) Shift() {
	o.Mux.Lock()
	headerHash := o.Order[0]
	delete(o.Map, headerHash)
	o.Order = o.Order[1:]
	o.Mux.Unlock()
}

// Len returns the OrderedWorkMap length
func (o *OrderedWorkMap) Len() int {
	o.Mux.Lock()
	out := len(o.Order)
	o.Mux.Unlock()
	return out
}

// WorkReceiver is a struct for the work receiver daemon
type WorkReceiver struct {
	httpServer        *http.Server
	shuttingDown      bool
	subscriptions     []chan []string
	subscriptionsMux  sync.Mutex
	lastWork          []string
	workHistory       OrderedWorkMap
	shareDiff         uint64
	shareTargetHex    string
	shareTargetBigInt *big.Int
	shareDiffBigInt   *big.Int
	BestShareTarget   *big.Int
}

// GetLastWork returns last work
func (w *WorkReceiver) GetLastWork(applyShareDiff bool) []string {
	work := w.lastWork
	// Apply Share Diff
	if applyShareDiff {
		work[2] = w.shareTargetHex
	}

	return work
}

// NewWorkReceiver creates new WorkReceiver instance
func NewWorkReceiver(bind string, shareDiff uint64) *WorkReceiver {
	shareTargetBigInt := big.NewInt(0).Div(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0)), big.NewInt(0).SetUint64(shareDiff))
	workReceiver := WorkReceiver{
		shareDiff:         shareDiff,
		shareDiffBigInt:   big.NewInt(0).SetUint64(shareDiff),
		shareTargetBigInt: shareTargetBigInt,
		shareTargetHex:    "0x" + hex.EncodeToString(utils.PadByteArrayStart(shareTargetBigInt.Bytes(), 32)),
		lastWork:          []string{"0x0", "0x0", "0x0", "0x0"},
		BestShareTarget:   big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0)),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			log.Logger.WithFields(logrus.Fields{
				"prefix":   "workreceiver",
				"expected": "POST",
				"got":      r.Method,
			}).Error("Invalid HTTP method")
			return
		}
		data, err := ioutil.ReadAll(r.Body)
		var parsedJSONData []string
		err = json.Unmarshal(data, &parsedJSONData)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"prefix": "workreceiver",
				"error":  err,
			}).Error("Unable to parse OpenEthereum work notification")
			return
		}

		if len(parsedJSONData) != 4 {
			log.Logger.WithFields(logrus.Fields{
				"prefix":   "workreceiver",
				"expected": 4,
				"got":      len(parsedJSONData),
			}).Error("Invalid work notification length (Ensure that you're using OpenEthereum)")
			return
		}

		var channelIndexesToClean []int

		workReceiver.lastWork = parsedJSONData

		workWithShareDifficulty := parsedJSONData
		workWithShareDifficulty[2] = workReceiver.shareTargetHex

		// Sending work notification to all subscribers
		workReceiver.subscriptionsMux.Lock()
		for i, ch := range workReceiver.subscriptions {
			if !isChanClosed(ch) {
				ch <- parsedJSONData
			} else {
				channelIndexesToClean = append(channelIndexesToClean, i)
			}
		}

		length := len(workReceiver.subscriptions)

		for _, chIndex := range channelIndexesToClean {
			workReceiver.subscriptions[chIndex] = workReceiver.subscriptions[length-1]
			workReceiver.subscriptions = workReceiver.subscriptions[:length-1]
		}
		workReceiver.subscriptionsMux.Unlock()

		workReceiver.workHistory.Append(parsedJSONData[0], parsedJSONData)

		if workReceiver.workHistory.Len() > 8 {
			// Removing unneeded (9th in history) work
			workReceiver.workHistory.Shift()
		}

		log.Logger.WithFields(logrus.Fields{
			"prefix":      "workreceiver",
			"header-hash": parsedJSONData[0][2:10],
		}).Info("New job for #" + strconv.FormatUint(utils.MustSoftHexToUint64(parsedJSONData[3]), 10))
	})

	workReceiver.httpServer = &http.Server{
		Addr:    bind,
		Handler: mux,
	}

	workReceiver.workHistory.Init()

	return &workReceiver
}

// Run function runs the WorkReceiver
func (w *WorkReceiver) Run() {
	err := w.httpServer.ListenAndServe()

	if !w.shuttingDown {
		panic(err)
	}
}

// Stop function stops the WorkReceiver
func (w *WorkReceiver) Stop() {
	err := w.httpServer.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

// SubscribeNotifications subscribes the given channel to the work receiver
func (w *WorkReceiver) SubscribeNotifications(ch chan []string) {
	w.subscriptions = append(w.subscriptions, ch)
}

func isChanClosed(ch <-chan []string) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}
