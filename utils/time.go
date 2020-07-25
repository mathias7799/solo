package utils

import "time"

// GetCurrent10MinTimestamp returns current timestamp without other 10 min remainder
func GetCurrent10MinTimestamp() int64 {
	return time.Now().Unix() / 600 * 600
}
