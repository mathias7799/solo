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
	"fmt"
	"strconv"

	"github.com/flexpool/solo/utils"
	"github.com/mitchellh/mapstructure"
)

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

// GetBlockByNumber delegates to `eth_getBlockByNumber` RPC method, and returns block by number
func (n *Node) GetBlockByNumber(blockNumber uint64) (Block, error) {
	data, err := n.makeHTTPRPCRequest("eth_getBlockByNumber", []interface{}{fmt.Sprintf("0x%x", blockNumber), false})
	if err != nil {
		return Block{}, err
	}

	var block Block
	err = mapstructure.Decode(data, &block)
	return block, err
}

// GetUncleByBlockNumberAndIndex delegates to `eth_getUncleByBlockNumberAndIndex` RPC method, and returns uncle by block number and index
func (n *Node) GetUncleByBlockNumberAndIndex(blockNumber uint64, uncleIndex int) (Block, error) {
	data, err := n.makeHTTPRPCRequest("eth_getUncleByBlockNumberAndIndex", []interface{}{fmt.Sprintf("0x%x", blockNumber), fmt.Sprintf("0x%x", uncleIndex)})
	if err != nil {
		return Block{}, err
	}

	var block Block
	err = mapstructure.Decode(data, &block)
	return block, err
}

// GetUncleCountByBlockNumber delegates to `eth_getUncleCountByBlockNumber` RPC method, and returns amount of uncles by given block number
func (n *Node) GetUncleCountByBlockNumber(blockNumber uint64) (uint64, error) {
	data, err := n.makeHTTPRPCRequest("eth_getUncleCountByBlockNumber", []interface{}{fmt.Sprintf("0x%x", blockNumber)})
	if err != nil {
		return 0, err
	}

	uncleCount, err := strconv.ParseUint(utils.Clear0x(data.(string)), 16, 64)
	if err != nil {
		return 0, err
	}

	return uncleCount, nil
}
