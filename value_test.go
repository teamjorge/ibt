package ibt

import (
	"reflect"
	"testing"

	"github.com/teamjorge/ibt/headers"
)

func TestValue(t *testing.T) {
	t.Run("test readVarValue byte", func(t *testing.T) {
		bitValue := []byte{0x0}
		bitArray := append(bitValue, bitValue...)

		varHeader := headers.VarHeader{
			Count:  1,
			Rtype:  0,
			Offset: 0,
		}

		value := readVarValue(bitValue, varHeader)

		if reflect.TypeOf(value).String() != "uint8" {
			t.Errorf("expected return value to be uint8. got %v of type %T", value, value)
		}

		if value.(uint8) != uint8(0) {
			t.Errorf("expected returned value to be %d. received %v", uint8(0), value.(uint8))
		}

		varHeader.Count = 2

		values := readVarValue(bitArray, varHeader)

		if reflect.TypeOf(values).String() != "[]uint8" {
			t.Errorf("expected return value to be []uint8. got %v of type %T", values, values)
		}

		if values.([]uint8)[0] != uint8(0) || values.([]uint8)[1] != uint8(0) {
			t.Errorf("expected returned values to be %v. received %v", []uint8{0, 0}, values)
		}
	})

	t.Run("test readVarValue bool", func(t *testing.T) {
		bitValue := []byte{0x1}
		bitArray := []byte{0x0, 0x1}

		varHeader := headers.VarHeader{
			Count:  1,
			Rtype:  1,
			Offset: 0,
		}

		value := readVarValue(bitValue, varHeader)

		if reflect.TypeOf(value).String() != "bool" {
			t.Errorf("expected return value to be boolean. got %v of type %T", value, value)
		}

		if value.(bool) != true {
			t.Errorf("expected returned value to be %v. received %v", true, value.(bool))
		}

		varHeader.Count = 2

		values := readVarValue(bitArray, varHeader)

		if reflect.TypeOf(values).String() != "[]bool" {
			t.Errorf("expected return value to be []bool. got %v of type %T", values, values)
		}

		if values.([]bool)[0] || !values.([]bool)[1] {
			t.Errorf("expected returned values to be %v. received %v", []bool{false, true}, values)
		}
	})

	t.Run("test readVarValue int", func(t *testing.T) {
		bitValue := []byte{0x1, 0x0, 0x0, 0x0}
		bitArray := []byte{0x1, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0}

		varHeader := headers.VarHeader{
			Count:  1,
			Rtype:  2,
			Offset: 0,
		}

		value := readVarValue(bitValue, varHeader)

		if reflect.TypeOf(value).String() != "int" {
			t.Errorf("expected return value to be int. got %v of type %T", value, value)
		}

		if value.(int) != 1 {
			t.Errorf("expected returned value to be %v. received %v", 1, value.(int))
		}

		varHeader.Count = 2

		values := readVarValue(bitArray, varHeader)

		if reflect.TypeOf(values).String() != "[]int" {
			t.Errorf("expected return value to be []int. got %v of type %T", values, values)
		}

		if values.([]int)[0] != 1 || values.([]int)[1] != 2 {
			t.Errorf("expected returned values to be %v. received %v", []int{1, 2}, values)
		}
	})

	t.Run("test readVarValue bitfield", func(t *testing.T) {
		bitValue := []byte{0x0, 0x2, 0x4, 0x10}
		bitArray := []byte{0x0, 0x2, 0x4, 0x10, 0x0, 0x4, 0x8, 0x14}

		varHeader := headers.VarHeader{
			Count:  1,
			Rtype:  3,
			Offset: 0,
		}

		value := readVarValue(bitValue, varHeader)

		if reflect.TypeOf(value).String() != "string" {
			t.Errorf("expected return value to be string. got %v of type %T", value, value)
		}

		if value.(string) != "0x10040200" {
			t.Errorf("expected returned value to be %v. received %v", "0x10040200", value.(string))
		}

		varHeader.Count = 2

		values := readVarValue(bitArray, varHeader)

		if reflect.TypeOf(values).String() != "[]string" {
			t.Errorf("expected return value to be []string. got %v of type %T", values, values)
		}

		if values.([]string)[0] != "0x10040200" || values.([]string)[1] != "0x14080400" {
			t.Errorf("expected returned values to be %v. received %v", []string{"0x10040200", "0x14080400"}, values)
		}
	})

	t.Run("test readVarValue float32", func(t *testing.T) {
		bitValue := []byte{0x1, 0x7c, 0x17, 0xba}
		bitArray := []byte{0x1, 0x7c, 0x17, 0xba, 0x1, 0x9c, 0x12, 0xba}

		varHeader := headers.VarHeader{
			Count:  1,
			Rtype:  4,
			Offset: 0,
		}

		value := readVarValue(bitValue, varHeader)

		if reflect.TypeOf(value).String() != "float32" {
			t.Errorf("expected return value to be float32. got %v of type %T", value, value)
		}

		if value.(float32) != -0.0005778671 {
			t.Errorf("expected returned value to be %v. received %v", -0.0005778671, value.(float32))
		}

		varHeader.Count = 2

		values := readVarValue(bitArray, varHeader)

		if reflect.TypeOf(values).String() != "[]float32" {
			t.Errorf("expected return value to be []float32. got %v of type %T", values, values)
		}

		if values.([]float32)[0] != -0.0005778671 || values.([]float32)[1] != -0.00055927044 {
			t.Errorf("expected returned values to be %v. received %v", []float32{-0.0005778671, -0.00055927044}, values)
		}
	})

	t.Run("test readVarValue float64", func(t *testing.T) {
		bitValue := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x75, 0x22, 0x41}
		bitArray := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x75, 0x22, 0x41, 0x1, 0x44, 0x0, 0x0, 0x0, 0x75, 0x22, 0x41}

		varHeader := headers.VarHeader{
			Count:  1,
			Rtype:  5,
			Offset: 0,
		}

		value := readVarValue(bitValue, varHeader)

		if reflect.TypeOf(value).String() != "float64" {
			t.Errorf("expected return value to be float64. got %v of type %T", value, value)
		}

		if value.(float64) != 604800.0 {
			t.Errorf("expected returned value to be %v. received %v", 604800.0, value.(float64))
		}

		varHeader.Count = 2

		values := readVarValue(bitArray, varHeader)

		if reflect.TypeOf(values).String() != "[]float64" {
			t.Errorf("expected return value to be []float64. got %v of type %T", values, values)
		}

		if values.([]float64)[0] != 604800 || values.([]float64)[1] != 604800.0000020267 {
			t.Errorf("expected returned values to be %v. received %v", []float64{604800, 604800.0000020267}, values)
		}
	})

	// one := []byte{0x0}
	// two := []byte{0x1, 0x0, 0x0, 0x0}
	// three := []byte{0x0, 0x2, 0x4, 0x10}
	// four := []byte{0x1, 0x7c, 0x17, 0xba}
	// five := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x75, 0x22, 0x41}

	// fourArr := []byte{0xc0, 0x42, 0x56, 0xc1, 0x8e, 0x7d, 0x30, 0xc2}
}
