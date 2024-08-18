package utilities

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

const invalidValue int = 4294967295

// Byte4ToInt will convert the little endian 4 byte value into an int
func Byte4ToInt(in []byte) int {
	converted := int(binary.LittleEndian.Uint32(in))
	if converted == invalidValue {
		return -1
	}

	return converted
}

// Byte4ToFloat will convert the little endian 4 byte value into an float
func Byte4ToFloat(in []byte) float32 {
	bits := binary.LittleEndian.Uint32(in)
	return math.Float32frombits(bits)
}

// Byte8ToFloat will convert the little endian 8 byte value into an float
func Byte8ToFloat(in []byte) float64 {
	bits := binary.LittleEndian.Uint64(in)
	return math.Float64frombits(bits)
}

// Byte4toBitField will convert the little endian 8 byte value into a bitfield string value
func Byte4toBitField(in []byte) string {
	return fmt.Sprintf("0x%x", int(binary.LittleEndian.Uint32(in)))
}

// BytesToString will convert the given bytes to string and remove any additional padding bytes
func BytesToString(in []byte) string {
	return strings.TrimRight(string(in), "\x00")
}

// Byte8ToInt will convert the little endian 8 byte value into an int
func Byte8ToInt64(in []byte) int64 {
	bits := binary.LittleEndian.Uint64(in)
	return int64(bits)
}
