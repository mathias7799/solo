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

package utils

import (
	"errors"
	"fmt"
	"math/big"
	"net"
	"strconv"
)

// Clear0x removes 0x prefix from hex string
func Clear0x(s string) string {
	if len(s) < 2 {
		fmt.Println("utils.clear0x: Unexpected input\"" + s + "\"")
		return ""
	}
	if s[:2] == "0x" {
		return s[2:]
	}
	return s
}

// MustSoftHexToUint64 is a must function that converts Hex to Uint64 without panic
func MustSoftHexToUint64(s string) uint64 {
	out, err := strconv.ParseUint(Clear0x(s), 16, 64)
	if err != nil {
		fmt.Println("utils.MustSoftHexToUint64L Invalid hex string \"" + s + "\"")
	}
	return out
}

// IsInvalidAddress checks if given address is invalid
func IsInvalidAddress(s string) error {
	host, _, err := net.SplitHostPort(s)
	if err != nil {
		return err
	}

	if net.ParseIP(host) == nil {
		return errors.New("Invalid IP address")
	}

	return nil
}

// HexStrToBigInt converts hex string to big.Int
func HexStrToBigInt(hexStr string) *big.Int {
	v := new(big.Int)
	v.SetString(Clear0x(hexStr), 16)
	if v == nil {
		return big.NewInt(0)
	}
	return v
}
