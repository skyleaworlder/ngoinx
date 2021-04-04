package utils

import (
	"crypto/sha1"
	"encoding/binary"
)

// GetHashID is a tool function
// and this function return HID (int)
func GetHashID(input string) (HID uint64) {
	buf := []byte(input)
	sha1Sum := sha1.Sum(buf)
	return binary.BigEndian.Uint64(sha1Sum[:][:8])
}
