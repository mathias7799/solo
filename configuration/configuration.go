package configuration

import "github.com/kelseyhightower/envconfig"

// Configuration specifies the Solo configuration
type Configuration struct {
	WorkmanagerNotificationsBindAddr string `envconfig:"solo_workrmanager_notifications_bind_addr" required:"true"`
	GatewayInsecureBindAddr          string `envconfig:"solo_gateway_insecure_bind_addr"`
	GatewayPassword                  string `envconfig:"solo_gateway_password" required:"true"`
	ShareDifficulty                  uint64 `envconfig:"solo_share_difficulty" default:"4000000000"`
	NodeHTTPRPC                      string `envconfig:"solo_node_http_rpc" default:"127.0.0.1:8545"`
}

// GetConfig parses the environment variables
// and returns them as Configuration struct
func GetConfig() (Configuration, error) {
	c := Configuration{}
	err := envconfig.Process("solo", &c)
	return c, err
}
