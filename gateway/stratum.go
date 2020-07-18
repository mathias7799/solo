package gateway

import (
	"bufio"
	"fmt"
	"net"

	"github.com/flexpool/solo/jsonrpc"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/types"
	"github.com/sirupsen/logrus"
)

func write(conn net.Conn, data []byte) (int, error) {
	return conn.Write(append(data, '\n'))
}

// RunWorkSender runs a work sender for a given connection
func (g *Gateway) RunWorkSender(conn net.Conn) {
	// Creating a channel and subscribing it to the work receiver
	ch := make(chan []string)

	g.parentWorkReceiver.SubscribeNotifications(ch)

	for {
		work := <-ch
		_, err := write(conn, jsonrpc.MarshalResponse(jsonrpc.Response{
			JSONRPCVersion: jsonrpc.Version,
			ID:             0,
			Result:         work,
		}))

		if err != nil {
			break
		}
	}

	// Closed channel would be automatically unsubscribed by work receiver garbege collector
	close(ch)
}

// HandleConnection handles the gateway connection
func (g *Gateway) HandleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	var authenticated bool

	var workerName string

	ip := conn.RemoteAddr().String()

	for scanner.Scan() {
		request, err := jsonrpc.UnmarshalRequest(scanner.Bytes())
		if err != nil {
			write(conn, GetInvalidRequestError(0))
			log.Logger.WithFields(logrus.Fields{
				"prefix": "gateway",
				"ip":     ip,
			}).Warn("Invalid JSONRPC request")

			// Close connection if not authenticated
			if !authenticated {
				return
			}
			continue
		}

		if !authenticated {
			if request.Method != "eth_submitLogin" {
				write(conn, GetUnauthorizedError(request.ID))
				return
			}

			if len(request.Params) < 2 {
				write(conn, GetInvalidCredentialsError(request.ID))
				return
			}

			workerName = request.Params[0]

			if request.Params[1] != g.stratumPassword {
				log.Logger.WithFields(logrus.Fields{
					"prefix":      "gateway",
					"worker-name": workerName,
					"ip":          ip,
				}).Warn("Invalid password")
				write(conn, GetInvalidCredentialsError(request.ID))
				return
			}

			log.Logger.WithFields(logrus.Fields{
				"prefix":      "gateway",
				"worker-name": workerName,
				"ip":          ip,
			}).Info("Authenticated new worker")

			write(conn, jsonrpc.MarshalResponse(jsonrpc.Response{
				JSONRPCVersion: jsonrpc.Version,
				ID:             request.ID,
				Result:         true,
				Error:          nil,
			}))

			authenticated = true

			// Starting work sender
			go g.RunWorkSender(conn)

			continue
		}

		switch request.Method {
		case "eth_getWork":
			write(conn, jsonrpc.MarshalResponse(jsonrpc.Response{
				JSONRPCVersion: jsonrpc.Version,
				ID:             request.ID,
				Result:         g.parentWorkReceiver.GetLastWork(true),
				Error:          nil,
			}))
		case "eth_submitWork":
			if len(request.Params) < 3 || len(request.Params[0]) != 18 || len(request.Params[1]) != 66 || len(request.Params[2]) != 66 {
				write(conn, GetInvalidParamsError(request.ID))
			} else {
				shareType, err := g.submitShare(request.Params, workerName)
				if err != nil {
					write(conn, jsonrpc.MarshalResponse(jsonrpc.Response{
						JSONRPCVersion: jsonrpc.Version,
						ID:             request.ID,
						Result:         nil,
						Error:          err.Error(),
					}))
					continue
				}

				if shareType == 0 || shareType == 1 {
					write(conn, jsonrpc.MarshalResponse(jsonrpc.Response{
						JSONRPCVersion: jsonrpc.Version,
						ID:             request.ID,
						Result:         true,
						Error:          nil,
					}))
				} else {
					write(conn, GetInvalidShareError(request.ID))
				}

				log.Logger.WithFields(logrus.Fields{
					"prefix": "gateway",
					"worker": workerName,
					"ip":     ip,
				}).Info("Received " + types.ShareTypeNameMap[shareType] + " share")
			}
		case "eth_submitHashrate":
			// Return true
			write(conn, jsonrpc.MarshalResponse(jsonrpc.Response{
				JSONRPCVersion: jsonrpc.Version,
				ID:             request.ID,
				Result:         true,
				Error:          nil,
			}))
			fmt.Println("eth_submitHashrate was called, but is not implemented")
		default:
			write(conn, jsonrpc.MarshalResponse(jsonrpc.Response{
				JSONRPCVersion: jsonrpc.Version,
				ID:             request.ID,
				Result:         nil,
				Error:          "Method not found",
			}))
		}
	}
}
