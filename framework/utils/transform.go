package utils

import (
	"encoding/binary"
	"math"
)

func ByteArrayToUint16(bytes []byte) uint16 {
	return binary.LittleEndian.Uint16(bytes)
}

func Uint16ToByteArray(value uint16) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, value)
	return bytes
}

func ByteArrayToUint32(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes)
}

func Uint32ToByteArray(value uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, value)
	return bytes
}

func Uint64ToByteArray(value uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, value)
	return bytes
}

func ByteArrayToUint64(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(bytes)
}

func ConcatTwoUint32ToByteArray(secret1 uint32, secret2 uint32) []byte {
	newKey := make([]byte, 8)
	secret1Bytes := make([]byte, 4)
	secret2Bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(secret1Bytes, secret1)
	binary.LittleEndian.PutUint32(secret2Bytes, secret2)
	for i := 0; i < 4; i++ {
		newKey[i] = secret1Bytes[i]
		newKey[i+4] = secret2Bytes[i]
	}
	return newKey
}

func Float32FromByteArray(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func Float64FromByteArray(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
