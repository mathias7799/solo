package jsonrpc

import "encoding/json"

// Version is a JSONRPC version
const Version = "2.0"

// RequestStringParams specifies the JSONRPC gateway request
type RequestStringParams struct {
	JSONRPCVersion string   `json:"jsonrpc"`
	ID             int      `json:"id"`
	Method         string   `json:"method"`
	Params         []string `json:"params"`
}

// Request specifies the JSONRPC gateway request
type Request struct {
	JSONRPCVersion string      `json:"jsonrpc"`
	ID             int         `json:"id"`
	Method         string      `json:"method"`
	Params         interface{} `json:"params"`
}

// Response specifies the JSONRPC gateway response
type Response struct {
	JSONRPCVersion string      `json:"jsonrpc"`
	ID             int         `json:"id"`
	Result         interface{} `json:"result"`
	Error          interface{} `json:"error"`
}

// UnmarshalRequest parses the JSONRPC request, and returns it as a Request struct
func UnmarshalRequest(b []byte) (RequestStringParams, error) {
	var req RequestStringParams
	err := json.Unmarshal(b, &req)
	return req, err
}

// UnmarshalResponse parses the JSONRPC request, and returns it as a Response struct
func UnmarshalResponse(b []byte) (Response, error) {
	var resp Response
	err := json.Unmarshal(b, &resp)
	return resp, err
}

// MarshalResponse creates a JSONRPC response bytes from a Response struct
func MarshalResponse(r Response) []byte {
	resp, _ := json.Marshal(r)
	return resp
}

// MarshalRequestStringParams creates a JSONRPC request bytes from a RequestStringParams struct
func MarshalRequestStringParams(r RequestStringParams) []byte {
	req, _ := json.Marshal(r)
	return req
}

// MarshalRequest JSONRPC request bytes from a Request struct
func MarshalRequest(r Request) []byte {
	req, _ := json.Marshal(r)
	return req
}
