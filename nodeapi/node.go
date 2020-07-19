package nodeapi

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/flexpool/solo/jsonrpc"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/types"
	"github.com/flexpool/solo/utils"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Node is the base OpenEthereum API struct
type Node struct {
	httpRPCEndpoint string
	Type            types.NodeType
}

// NewNode creates a new Node instance
func NewNode(httpRPCEndpoint string) (*Node, error) {
	if _, err := url.Parse(httpRPCEndpoint); err != nil {
		return nil, errors.New("invalid HTTP URL")
	}

	node := Node{httpRPCEndpoint: httpRPCEndpoint}

	// Detecting node type
	clientVersion, err := node.ClientVersion()
	if err != nil {
		return nil, errors.Wrap(err, "failed detecting node type")
	}

	clientVersionSplitted := strings.Split(clientVersion, "/")
	if len(clientVersionSplitted) < 1 {
		log.Logger.WithFields(logrus.Fields{
			"prefix": "node",
		}).Warn("Unable to detect node type: received invalid client version. Falling back to Geth.")
	} else {
		clientVersion = clientVersionSplitted[0]
		switch strings.ToLower(clientVersion) {
		case "geth":
			node.Type = types.GethNode
		case "openethereum":
			node.Type = types.OpenEthereumNode
		default:
			log.Logger.WithFields(logrus.Fields{
				"prefix": "node",
			}).Warn("Unknown node \"" + clientVersion + "\". Falling back to Geth.")
		}
	}

	log.Logger.WithFields(logrus.Fields{
		"prefix": "node",
	}).Info("Configured for " + types.NodeStringMap[node.Type] + " node")

	return &node, nil
}

func (n *Node) makeHTTPRPCRequest(method string, params []string) (interface{}, error) {
	req := jsonrpc.MarshalRequest(jsonrpc.Request{
		JSONRPCVersion: jsonrpc.Version,
		ID:             rand.Intn(99999999),
		Method:         method,
		Params:         params,
	})

	response, err := http.Post(n.httpRPCEndpoint, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Additional error check
	parsedData, err := jsonrpc.UnmarshalResponse(data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal node's response ("+string(data)+")")
	}

	if parsedData.Error != nil {
		return nil, errors.New("unexpected node response: " + string(data))
	}

	return parsedData.Result, nil
}

// SubmitWork delegates to `eth_submitWork` API method, and submits work
func (n *Node) SubmitWork(work []string) (bool, error) {
	data, err := n.makeHTTPRPCRequest("eth_submitWork", work)
	if err != nil {
		return false, err
	}

	return data.(bool), nil
}

// BlockNumber delegates to `eth_blockNumber` API method, and returns the current block number
func (n *Node) BlockNumber() (uint64, error) {
	data, err := n.makeHTTPRPCRequest("eth_blockNumber", nil)
	if err != nil {
		return 0, err
	}

	blockNumber, err := strconv.ParseUint(utils.Clear0x(data.(string)), 16, 64)
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

// ClientVersion delegates to `eth_blockNumber` API method, and returns the current block number
func (n *Node) ClientVersion() (string, error) {
	data, err := n.makeHTTPRPCRequest("web3_clientVersion", nil)
	if err != nil {
		return "", err
	}

	return data.(string), nil
}
