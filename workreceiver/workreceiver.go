package workreceiver

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/utils"

	"github.com/sirupsen/logrus"
)

// WorkReceiver is a struct for the work receiver daemon
type WorkReceiver struct {
	httpServer   *http.Server
	shuttingDown bool
}

// NewWorkReceiver creates new WorkReceiver instance
func NewWorkReceiver(bind string) WorkReceiver {
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

		log.Logger.WithFields(logrus.Fields{
			"prefix":      "workreceiver",
			"header-hash": parsedJSONData[0][2:10],
		}).Info("New job for #" + strconv.FormatUint(utils.MustSoftHexToUint64(parsedJSONData[3]), 10))

	})

	return WorkReceiver{httpServer: &http.Server{
		Addr:    bind,
		Handler: mux,
	}}
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
