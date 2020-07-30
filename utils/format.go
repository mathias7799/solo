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

import "math"

// GetSI calculates the si prefix of the given number, and returns siDiv and siChar
func GetSI(number float64) (float64, string) {
	if number < 1000 {
		return 1, ""
	}
	symbols := "kMGTPEZY"
	symbolsLen := len(symbols)
	i := 1
	for {
		number /= 1000
		if number < 1000 || i == symbolsLen-1 {
			return math.Pow(1000, float64(i)), string(symbols[i-1])
		}
		i++
	}
}
