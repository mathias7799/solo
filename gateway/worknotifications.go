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

package gateway

import "encoding/json"

// OpenEthereumWorkNotification represents the struct of OpenEthereum work notification
type OpenEthereumWorkNotification struct {
	Result []string `json:"result"`
}

func parseGethWorkNotification(data []byte) ([]string, error) {
	var parsedData []string
	err := json.Unmarshal(data, &parsedData)
	if err != nil {
		return nil, err
	}
	return parsedData, nil
}

func parseOpenEthereumWorkNotification(data []byte) ([]string, error) {
	var parsedData OpenEthereumWorkNotification
	err := json.Unmarshal(data, &parsedData)
	if err != nil {
		return nil, err
	}
	return parsedData.Result, nil
}
