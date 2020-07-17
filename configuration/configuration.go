package configuration

import "github.com/kelseyhightower/envconfig"

// Configuration specifies the Solo configuration
type Configuration struct {
	WorkreceiverBindAddr string `envconfig:"solo_workreceiver_bind_addr" required:"true"`
}

// GetConfig parses the environment variables
// and returns them as Configuration struct
func GetConfig() (Configuration, error) {
	c := Configuration{}
	err := envconfig.Process("solo", &c)
	return c, err
}
