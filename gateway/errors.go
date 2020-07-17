package gateway

import "github.com/flexpool/solo/jsonrpc"

func marshalError(id int, msg string, result interface{}) []byte {
	return jsonrpc.MarshalResponse(jsonrpc.Response{
		JSONRPCVersion: jsonrpc.Version,
		ID:             id,
		Result:         result,
		Error:          msg,
	})
}

// GetInvalidRequestError creates and returns Stratum `Invalid JSONRPC Request` message
func GetInvalidRequestError(id int) []byte {
	return marshalError(id, "Invalid JSONRPC Request", nil)
}

// GetUnauthorizedError creates and restart Stratum `Unauthorized` message
func GetUnauthorizedError(id int) []byte {
	return marshalError(id, "Unauthorized", nil)
}

// GetInvalidParamsError creates and restart Stratum `Invalid parameters` message
func GetInvalidParamsError(id int) []byte {
	return marshalError(id, "Invalid Parameters", nil)
}

// GetInvalidCredentialsError creates and restart Stratum `Invalid credentials` message
func GetInvalidCredentialsError(id int) []byte {
	return marshalError(id, "Invalid credentials", nil)
}

// GetNotRequestedWorkError creates and restart Stratum `Work is outdated, or not requested` message
func GetNotRequestedWorkError(id int) []byte {
	return marshalError(id, "Work is outdated, or not requested", nil)
}

// GetInvalidShareError creates and restart Stratum `Provided POW solution is invalid` message
func GetInvalidShareError(id int) []byte {
	return marshalError(id, "Provided POW solution is invalid", false)
}
