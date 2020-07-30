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

package types

// NodeType is a shortcut for uint64, and is used to store information about which node it is
type NodeType uint64

// GethNode specifies NodeType for Geth (go-ethereum)
const GethNode = 0

// OpenEthereumNode specifies NodeType for OpenEthereum (formerly Parity Ethereum)
const OpenEthereumNode = 1

// NodeStringMap is a mapping that allows to quickly access the name of node
var NodeStringMap = map[NodeType]string{
	GethNode:         "Geth",
	OpenEthereumNode: "OpenEthereum",
}
