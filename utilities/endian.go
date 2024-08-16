package utilities

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

const invalidValue int = 4294967295

func Byte4ToInt(in []byte) int {
	converted := int(binary.LittleEndian.Uint32(in))
	if converted == invalidValue {
		return -1
	}

	return converted
}

func Byte4ToFloat(in []byte) float32 {
	bits := binary.LittleEndian.Uint32(in)
	return math.Float32frombits(bits)
}

func Byte8ToFloat(in []byte) float64 {
	bits := binary.LittleEndian.Uint64(in)
	return math.Float64frombits(bits)
}

func Byte4toBitField(in []byte) string {
	return fmt.Sprintf("0x%x", int(binary.LittleEndian.Uint32(in)))
}

func BytesToString(in []byte) string {
	return strings.TrimRight(string(in), "\x00")
}

func Byte8ToInt64(in []byte) int64 {
	bits := binary.LittleEndian.Uint64(in)
	return int64(bits)
}
