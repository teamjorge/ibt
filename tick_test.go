package ibt

import "testing"

func TestGetTickValue(t *testing.T) {
	testTick := Tick{
		"Speed": float32(103.23),
		"Gear":  5,
		"Flag":  "0x3421",
	}

	t.Run("test normal scenario", func(t *testing.T) {
		value, err := GetTickValue[int](testTick, "Gear")
		if err != nil {
			t.Errorf("expected err to be nil but received: %v", err)
		}

		if value != 5 {
			t.Errorf("expected return Gear value to be %d. received: %d", 5, value)
		}
	})

	t.Run("test missing key", func(t *testing.T) {
		_, err := GetTickValue[int](testTick, "NotFound")
		if err == nil {
			t.Errorf("expected an error to occur when retrieving value for key %s", "NotFound")
		}
	})

	t.Run("test missing key", func(t *testing.T) {
		_, err := GetTickValue[int](testTick, "Speed")
		if err == nil {
			t.Errorf("expected an error to occur when retrieving value for key %s with type int", "Speed")
		}
	})
}
