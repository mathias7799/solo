package gateway

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"time"

	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/process"
	"github.com/flexpool/solo/utils"
	"github.com/sirupsen/logrus"
)

// Gateway is a stratum proxy that servers workers
type Gateway struct {
	bind              string
	stratumPassword   string
	isSecure          bool
	tlsKeyPair        tls.Certificate
	context           context.Context
	cancelContextFunc context.CancelFunc
}

// NewGatewayInsecure creates Non SSL gateway instance
func NewGatewayInsecure(bind string, password string) (Gateway, error) {
	err := utils.IsInvalidAddress(bind)
	if err != nil {
		return Gateway{}, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	return Gateway{bind: bind, stratumPassword: password, isSecure: false, context: ctx, cancelContextFunc: cancelFunc}, nil
}

// Run runs the Gateway
func (g *Gateway) Run() {
	laddr, err := net.ResolveTCPAddr("tcp", g.bind)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "gateway",
			"bind":   g.bind,
			"error":  err,
			"secure": g.isSecure,
		}).Error("Unable to resolve TCP Address")
		process.SafeExit(1)
		return
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "gateway",
			"bind":   g.bind,
			"error":  err,
			"secure": g.isSecure,
		}).Error("Unable to listen gateway port")
		process.SafeExit(1)
		return
	}

	log.Logger.WithFields(logrus.Fields{
		"prefix": "gateway",
		"bind":   g.bind,
		"secure": g.isSecure,
	}).Info("Started gateway")

	for {

		select {
		case <-g.context.Done():
			listener.Close()
			log.Logger.WithFields(logrus.Fields{
				"prefix": "gateway",
				"secure": g.isSecure,
			}).Info("Stopped server")
		default:
			listener.SetDeadline(time.Now().Add(time.Second))

			conn, err := listener.Accept()
			if err != nil {
				if os.IsTimeout(err) {
					continue
				}

				log.Logger.WithFields(logrus.Fields{
					"prefix": "gateway",
					"error":  err,
					"secure": g.isSecure,
				}).Error("Unable to accept TCP connection")
				continue
			}
			go g.HandleConnection(conn)
		}
	}
}

// Stop stops the gateway server
func (g *Gateway) Stop() {
	g.cancelContextFunc()
}
