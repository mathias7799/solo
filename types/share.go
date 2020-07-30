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

// ShareType is a shortcut for uint8
type ShareType uint8

const (
	// ShareValid is a shortcut for uint8(0)
	ShareValid = 0
	// ShareStale is a shortcut for uint8(1)
	ShareStale = 1
	// ShareInvalid is a shortcut for uint8(2)
	ShareInvalid = 2
)

// ShareTypeNameMap has ShareType => <Share name (String)> mapping
var ShareTypeNameMap = map[ShareType]string{
	0: "valid",
	1: "stale",
	2: "invalid",
}
