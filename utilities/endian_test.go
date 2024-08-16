package utilities

import (
	"testing"
)

func TestEndian(t *testing.T) {
	t.Run("test Byte4ToInt valid", func(t *testing.T) {
		res := Byte4ToInt([]byte{0x2, 0x0, 0x0, 0x1})

		if res != 16777218 {
			t.Errorf("expected result to be %d, got %d", 16777218, res)
		}
	})

	t.Run("test Byte4ToInt invalid", func(t *testing.T) {
		res := Byte4ToInt([]byte{0xff, 0xff, 0xff, 0xff})

		if res != -1 {
			t.Errorf("expected result to be %d, got %d", 0, res)
		}
	})

	t.Run("test Byte4ToFloat", func(t *testing.T) {
		res := Byte4ToFloat([]byte{0x1, 0x7c, 0x17, 0xba})

		if res != -0.0005778671 {
			t.Errorf("expected result to be %f, got %f", -0.0005778671, res)
		}
	})

	t.Run("test Byte8ToFloat valid", func(t *testing.T) {
		res := Byte8ToFloat([]byte{0x1, 0x44, 0x0, 0x0, 0x0, 0x75, 0x22, 0x41})

		if res != 604800.0000020267 {
			t.Errorf("expected result to be %f, got %f", 604800.0000020267, res)
		}
	})

	t.Run("test Byte4toBitField valid", func(t *testing.T) {
		res := Byte4toBitField([]byte{0x1, 0x44, 0x0, 0x0, 0x0, 0x75, 0x22, 0x41})

		if res != "0x4401" {
			t.Errorf("expected result to be %s, got %s", "0x4401", res)
		}
	})

	t.Run("test BytesToString valid", func(t *testing.T) {
		testStringNoSpace := []byte("testring")
		testStringWithSpace := append(testStringNoSpace, 0x00, 0x00)

		resNoSpace := BytesToString(testStringNoSpace)
		resWithSpace := BytesToString(testStringWithSpace)

		if resNoSpace != "testring" {
			t.Errorf("expected result to be %s, got %s", "testring", resNoSpace)
		}

		if resWithSpace != "testring" {
			t.Errorf("expected result to be %s, got %s", "testring", resWithSpace)
		}
	})

	t.Run("test Byte8ToInt64 valid", func(t *testing.T) {
		res := Byte8ToInt64([]byte{0x1, 0x44, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})

		if res != 17409 {
			t.Errorf("expected result to be %d, got %d", 17409, res)
		}
	})
}
