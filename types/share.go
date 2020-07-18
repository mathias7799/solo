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
