// Flexpool Solo - A lightweight SOLO Ethereum mining pool
// Copyright (C) 2020  Flexpool
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package nodeapi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/flexpool/solo/jsonrpc"
	"github.com/flexpool/solo/log"
	"github.com/flexpool/solo/types"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Node is the base OpenEthereum API struct
type Node struct {
	httpRPCEndpoint string
	Type            types.NodeType
}

// Block is a block body representation
type Block struct {
	Number       string   `json:"number"`
	Hash         string   `json:"hash"`
	Nonce        string   `json:"nonce"`
	Difficulty   string   `json:"difficulty"`
	Timestamp    string   `json:"timestamp"`
	Transactions []string `json:"transactions"`
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

func (n *Node) makeHTTPRPCRequest(method string, params interface{}) (interface{}, error) {
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

// HarvestBlockByNonce is the most compatible way to find out the block hash
func (n *Node) HarvestBlockByNonce(givenNonceHex string, givenNumber uint64) (block Block, uncleParent uint64, err error) {
	maxRecursion := 10
	blocksPerLoop := 100

	givenNumberHex := fmt.Sprintf("0x%x", givenNumber)

	for i := 0; i <= maxRecursion; i++ {
		currentBlockNumber, err := n.BlockNumber()
		if err != nil {
			panic(err)
		}

		for i := 0; i <= blocksPerLoop; i++ {
			blockNumber := currentBlockNumber - uint64(i)
			block, _ := n.GetBlockByNumber(blockNumber)
			if block.Nonce == givenNonceHex {
				return block, 0, nil
			}

			uncleCount, _ := n.GetUncleCountByBlockNumber(blockNumber)
			uncleCountInt := int(uncleCount)

			var uncle Block

			for i := 0; i < uncleCountInt; i++ {
				uncle, _ = n.GetUncleByBlockNumberAndIndex(blockNumber, i)
				if uncle.Nonce == givenNonceHex {
					return uncle, blockNumber, nil
				}
			}
		}
		time.Sleep(time.Second * 5)
	}

	return Block{}, 0, errors.New("unable to harvest block by hash (nonce: " + givenNonceHex + ", number:" + givenNumberHex + ")")
}
