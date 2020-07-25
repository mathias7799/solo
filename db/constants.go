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

package db

// StatPrefix is used to map stat db objects
const StatPrefix = "stat__"

// TotalStatPrefix is used to map summarized stat db objects
const TotalStatPrefix = "total__"

// BestSharePrefix is used to map best share db objects
const BestSharePrefix = "best__"

// MinedBlockPrefix is used to map block db objects
const MinedBlockPrefix = "block__"

// MinedValidSharesKey is used to identify if key is mined valid shares counter
const MinedValidSharesKey = "mined_valid_shares"

// AverageTotalHashrateKey is used to identify if key is average total hashrate item
const AverageTotalHashrateKey = "average_total_hashrate"

// TotalSharesKey is used to identify if key is total shares item
const TotalSharesKey = "total_shares"
