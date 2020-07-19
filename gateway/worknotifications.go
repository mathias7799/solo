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
