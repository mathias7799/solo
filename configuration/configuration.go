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

package configuration

import "github.com/kelseyhightower/envconfig"

// Configuration specifies the Solo configuration
type Configuration struct {
	WorkmanagerNotificationsBindAddr string `envconfig:"solo_workmanager_notifications_bind_addr" required:"true"`
	GatewayInsecureBindAddr          string `envconfig:"solo_gateway_insecure_bind_addr"`
	GatewayPassword                  string `envconfig:"solo_gateway_password" required:"true"`
	ShareDifficulty                  uint64 `envconfig:"solo_share_difficulty" default:"4000000000"`
	NodeHTTPRPC                      string `envconfig:"solo_node_http_rpc" default:"http://127.0.0.1:8545"`
	DBPath                           string `envconfig:"solo_db_path" default:"./solo_db"`
	LogLevel                         string `envconfig:"solo_log_level" default:"info"`
	BlockConfirmationsRequired       uint64 `envconfig:"solo_block_confirmations_required" default:"60"`
}

// GetConfig parses the environment variables
// and returns them as Configuration struct
func GetConfig() (Configuration, error) {
	c := Configuration{}
	err := envconfig.Process("solo", &c)
	return c, err
}
