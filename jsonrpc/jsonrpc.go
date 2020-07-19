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
