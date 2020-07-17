package workreceiver

import (
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
	bind      string
	httpMux   *http.ServeMux
	Exited    bool
	ExitError error
}

// NewWorkReceiver creates new WorkReceiver instance
func NewWorkReceiver(bind string) WorkReceiver {
	w := WorkReceiver{
		bind:    bind,
		httpMux: http.NewServeMux(),
	}

	w.httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	return w
}

// Start function starts the WorkReceiver
func (w *WorkReceiver) Start() {
	w.Exited = false

	// Write server error immediately after exit
	w.ExitError = http.ListenAndServe(w.bind, w.httpMux)

	w.Exited = true
}
