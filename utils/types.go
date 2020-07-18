package utils

import "math/big"

// BigMax256bit represents 2^256
var BigMax256bit = big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))
