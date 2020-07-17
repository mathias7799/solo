package utils

// PadByteArrayStart pads big endian integer to required amount of bytes
func PadByteArrayStart(bytearray []byte, totalBytes int) []byte {
	if len(bytearray) >= totalBytes {
		return bytearray
	}

	bytesToPad := totalBytes - len(bytearray)
	outBytes := make([]byte, totalBytes)

	for i := bytesToPad; i < totalBytes; i++ {
		outBytes[i] = bytearray[i-bytesToPad]
	}

	return outBytes
}
