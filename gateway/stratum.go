package gateway

import (
	"bufio"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/flexpool/solo/jsonrpc"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/types"
	"github.com/flexpool/solo/utils"
	"github.com/sirupsen/logrus"
)

func write(conn net.Conn, data []byte) (int, error) {
	return conn.Write(append(data, '\n'))
}

// RunWorkSender runs a work sender for a given connection
func (g *Gateway) RunWorkSender(conn net.Conn) {
	// Creating a channel and subscribing it to the work receiver
	ch := make(chan []string)

	g.parentWorkManager.SubscribeNotifications(ch)

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

	// Add 5 sec timeout
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	scanner := bufio.NewScanner(conn)

	var authenticated bool

	var workerName string

	ipSplitted := strings.Split(conn.RemoteAddr().String(), ":")
	ip := strings.Join(ipSplitted[:len(ipSplitted)-1], ":")

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

			g.statsCollector.Mux.Lock()
			pendingStat := g.statsCollector.PendingStats[workerName]
			pendingStat.IPAddress = ip
			g.statsCollector.PendingStats[workerName] = pendingStat
			g.statsCollector.Mux.Unlock()

			write(conn, jsonrpc.MarshalResponse(jsonrpc.Response{
				JSONRPCVersion: jsonrpc.Version,
				ID:             request.ID,
				Result:         true,
				Error:          nil,
			}))

			authenticated = true

			// Remove timeout
			conn.SetReadDeadline(time.Time{})

			// Starting work sender
			go g.RunWorkSender(conn)

			continue
		}

		switch request.Method {
		case "eth_getWork":
			write(conn, jsonrpc.MarshalResponse(jsonrpc.Response{
				JSONRPCVersion: jsonrpc.Version,
				ID:             request.ID,
				Result:         g.parentWorkManager.GetLastWork(true),
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

				g.statsCollector.Mux.Lock()
				pendingStat := g.statsCollector.PendingStats[workerName]

				switch shareType {
				case types.ShareValid:
					pendingStat.ValidShares++
					if g.statsCollector.Database.IncrValidShares() != nil {
						log.Logger.Error("Unable to increment valid shares counter")
					}
				case types.ShareStale:
					pendingStat.StaleShares++
				case types.ShareInvalid:
					pendingStat.InvalidShares++
				}

				g.statsCollector.PendingStats[workerName] = pendingStat
				g.statsCollector.Mux.Unlock()

				log.Logger.WithFields(logrus.Fields{
					"prefix": "gateway",
					"worker": workerName,
					"ip":     ip,
				}).Info("Received " + types.ShareTypeNameMap[shareType] + " share")
			}
		case "eth_submitHashrate":
			if len(request.Params) < 1 {
				write(conn, GetInvalidParamsError(request.ID))
				continue
			}

			reportedHashrateBigInt := utils.HexStrToBigInt(request.Params[0])

			g.statsCollector.Mux.Lock()
			pendingStat := g.statsCollector.PendingStats[workerName]
			pendingStat.ReportedHashrate, _ = big.NewFloat(0).SetInt(reportedHashrateBigInt).Float64()
			g.statsCollector.PendingStats[workerName] = pendingStat
			g.statsCollector.Mux.Unlock()

			write(conn, jsonrpc.MarshalResponse(jsonrpc.Response{
				JSONRPCVersion: jsonrpc.Version,
				ID:             request.ID,
				Result:         true,
				Error:          nil,
			}))

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
