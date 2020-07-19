package gateway

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/nodeapi"
	"github.com/flexpool/solo/types"
	"github.com/flexpool/solo/utils"
	"github.com/pkg/errors"

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

// WorkManager is a struct for the work manager daemon
type WorkManager struct {
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
	Node              *nodeapi.Node
	engineWaitGroup   *sync.WaitGroup
}

// GetLastWork returns last work
func (w *WorkManager) GetLastWork(applyShareDiff bool) []string {
	work := w.lastWork
	// Apply Share Diff
	if applyShareDiff {
		work[2] = w.shareTargetHex
	}

	return work
}

// NewWorkManager creates new WorkManager instance
func NewWorkManager(bind string, shareDiff uint64, node *nodeapi.Node, engineWaitGroup *sync.WaitGroup) *WorkManager {
	shareTargetBigInt := big.NewInt(0).Div(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0)), big.NewInt(0).SetUint64(shareDiff))
	workManager := WorkManager{
		shareDiff:         shareDiff,
		shareDiffBigInt:   big.NewInt(0).SetUint64(shareDiff),
		shareTargetBigInt: shareTargetBigInt,
		shareTargetHex:    "0x" + hex.EncodeToString(utils.PadByteArrayStart(shareTargetBigInt.Bytes(), 32)),
		lastWork:          []string{"0x0", "0x0", "0x0", "0x0"},
		BestShareTarget:   big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0)),
		Node:              node,
		engineWaitGroup:   engineWaitGroup,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			log.Logger.WithFields(logrus.Fields{
				"prefix":   "workmanager",
				"expected": "POST",
				"got":      r.Method,
			}).Error("Invalid HTTP method")
			return
		}
		data, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"prefix":    "workmanager",
				"error":     err,
				"node-type": types.NodeStringMap[workManager.Node.Type],
			}).Error("Unable to read work notification")
			return
		}

		var workNotification []string
		var workNotificationParseError error

		switch workManager.Node.Type {
		case types.GethNode:
			workNotification, workNotificationParseError = parseGethWorkNotification(data)
		case types.OpenEthereumNode:
			workNotification, workNotificationParseError = parseOpenEthereumWorkNotification(data)
		default:
			panic("unknown node type " + strconv.Itoa(types.OpenEthereumNode))
		}

		if workNotificationParseError != nil {
			log.Logger.WithFields(logrus.Fields{
				"prefix":    "workmanager",
				"error":     workNotificationParseError,
				"node-type": types.NodeStringMap[workManager.Node.Type],
			}).Error("Unable to parse work notification")
			return
		}

		var channelIndexesToClean []int

		workManager.lastWork = workNotification

		workWithShareDifficulty := make([]string, 4)
		copy(workWithShareDifficulty, workNotification)
		workWithShareDifficulty[2] = workManager.shareTargetHex

		// Sending work notification to all subscribers
		workManager.subscriptionsMux.Lock()
		for i, ch := range workManager.subscriptions {
			if !isChanClosed(ch) {
				ch <- workWithShareDifficulty
			} else {
				channelIndexesToClean = append(channelIndexesToClean, i)
			}
		}

		length := len(workManager.subscriptions)

		for _, chIndex := range channelIndexesToClean {
			workManager.subscriptions[chIndex] = workManager.subscriptions[length-1]
			workManager.subscriptions = workManager.subscriptions[:length-1]
		}
		workManager.subscriptionsMux.Unlock()
		workManager.workHistory.Append(workNotification[0], workNotification)

		if workManager.workHistory.Len() > 8 {
			// Removing unneeded (9th in history) work
			workManager.workHistory.Shift()
		}

		workTarget, _ := big.NewInt(0).SetString(utils.Clear0x(workNotification[2]), 16)
		workDifficulty, _ := big.NewFloat(0).SetInt(big.NewInt(0).Div(utils.BigMax256bit, workTarget)).Float64()

		log.Logger.WithFields(logrus.Fields{
			"prefix":      "workmanager",
			"header-hash": workNotification[0][2:10],
			"block-diff":  humanize.SIWithDigits(workDifficulty, 2, "H"),
		}).Info("New job for #" + strconv.FormatUint(utils.MustSoftHexToUint64(workNotification[3]), 10))
	})

	workManager.httpServer = &http.Server{
		Addr:    bind,
		Handler: mux,
	}

	workManager.workHistory.Init()

	return &workManager
}

// Run function runs the WorkManager
func (w *WorkManager) Run() {
	w.engineWaitGroup.Add(1)
	err := w.httpServer.ListenAndServe()

	if !w.shuttingDown {
		panic(errors.Wrap(err, "Server shut down unexpectedly"))
	}
}

// Stop function stops the WorkManager
func (w *WorkManager) Stop() {
	w.shuttingDown = true
	err := w.httpServer.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
	w.engineWaitGroup.Done()
}

// SubscribeNotifications subscribes the given channel to the work receiver
func (w *WorkManager) SubscribeNotifications(ch chan []string) {
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
