package utils

import (
	"fmt"
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

func MustSoftHexToUint64(s string) uint64 {
	out, err := strconv.ParseUint(Clear0x(s), 16, 64)
	if err != nil {
		fmt.Println("utils.MustSoftHexToUint64L Invalid hex string \"" + s + "\"")
	}
	return out
}
