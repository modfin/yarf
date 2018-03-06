package tnats

import "encoding/binary"

func intToBytes(i int) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(i))
	return bytes
}

func bytesToInt(b []byte) int {
	return int(binary.LittleEndian.Uint32(b))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
