package nodeapi

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/flexpool/solo/jsonrpc"
	"github.com/pkg/errors"
)

// Node is the base OpenEthereum API struct
type Node struct {
	httpRPCEndpoint string
}

// NewNode creates a new Node instance
func NewNode(httpRPCEndpoint string) (*Node, error) {
	if _, err := url.Parse(httpRPCEndpoint); err != nil {
		return nil, errors.Wrap(err, "Invalid HTTP URL")
	}

	return &Node{httpRPCEndpoint: httpRPCEndpoint}, nil
}

func (n *Node) makeHTTPRPCRequest() ([]byte, error) {
	req := jsonrpc.MarshalRequest(jsonrpc.Request{
		JSONRPCVersion: jsonrpc.Version,
		ID:             rand.Intn(99999999),
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

	// Additional error unmarshaling

	return data, nil
}

// SubmitWork delegates to `eth_submitWork` API method, and submits work
func (n *Node) SubmitWork(work []string) {

}
