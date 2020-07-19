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
