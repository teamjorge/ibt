package ibt

import (
	"bytes"
	"crypto/rand"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/teamjorge/ibt/headers"
)

func TestParser(t *testing.T) {
	f, err := os.Open(".testing/valid_test_file.ibt")
	if err != nil {
		t.Errorf("failed to open testing file - %v", err)
		return
	}
	defer f.Close()

	testHeaders, err := headers.ParseHeaders(f)
	if err != nil {
		t.Errorf("failed to parse header for testing file - %v", err)
		return
	}

	t.Run("test NewParser", func(t *testing.T) {
		p := NewParser(f, testHeaders, "Speed", "Lap")

		if p.bufferOffset != 53764 {
			t.Errorf("expected bufferOffset to be %d, recieved: %d", 53764, p.bufferOffset)
		}

		if p.length != 1072 {
			t.Errorf("expected length to be %d, recieved: %d", 1072, p.length)
		}

		expectedWhitelist := []string{"Speed", "Lap"}

		if p.whitelist[0] != expectedWhitelist[0] || p.whitelist[1] != expectedWhitelist[1] {
			t.Errorf("expected whitelist to be %v, recieved: %v", expectedWhitelist, p.whitelist)
		}

		if len(p.varHeader) != 276 {
			t.Errorf("expected varHeader to be of length %d, actual: %d", 276, len(p.varHeader))
		}
	})

	t.Run("test NewParser wildcard or null whitelist", func(t *testing.T) {
		p := NewParser(f, testHeaders, "*")

		if len(p.whitelist) != 276 {
			t.Errorf("expected whitelist to be of length %d, actual: %d", 278, len(p.whitelist))
		}

		p = NewParser(f, testHeaders)
		if len(p.whitelist) != 276 {
			t.Errorf("expected whitelist to be of length %d, actual: %d", 278, len(p.whitelist))
		}
	})
}

func TestParserNext(t *testing.T) {
	f, err := os.Open(".testing/valid_test_file.ibt")
	if err != nil {
		t.Errorf("failed to open testing file - %v", err)
		return
	}
	defer f.Close()

	testHeaders, err := headers.ParseHeaders(f)
	if err != nil {
		t.Errorf("failed to parse header for testing file - %v", err)
		return
	}

	t.Run("test parser Next() normal", func(t *testing.T) {
		p := NewParser(f, testHeaders, "LapCurrentLapTime")

		expectedValues := []float32{
			37.678566,
			37.695232,
			37.7119,
		}

		for idx, expectedValue := range expectedValues {
			vars, next := p.Next()
			if vars["LapCurrentLapTime"].Value != expectedValue {
				t.Errorf("expected LapCurrentLapTime value to equal %f, got %f", expectedValue, vars["LapCurrentLapTime"].Value)
			}
			if !next {
				t.Errorf("expected additional var values to be available after iteration %d", idx)
			}
		}
	})

	t.Run("test parser Next() reach end of buffer", func(t *testing.T) {
		p := NewParser(f, testHeaders, "LapCurrentLapTime")

		expectedValue1 := float32(44.128567)

		p.current = 388
		vars, next := p.Next()
		if vars["LapCurrentLapTime"].Value != expectedValue1 {
			t.Errorf("expected LapCurrentLapTime value to equal %f, got %f", expectedValue1, vars["LapCurrentLapTime"].Value)
		}
		if !next {
			t.Error("expected additional var values to be available after iteration")
		}

		expectedValue2 := float32(44.145233)
		vars, next = p.Next()
		if vars["LapCurrentLapTime"].Value != expectedValue2 {
			t.Errorf("expected LapCurrentLapTime value to equal %f, got %f", expectedValue2, vars["LapCurrentLapTime"].Value)
		}
		if next {
			t.Error("expected no more var values to be available after iteration")
		}
	})

	t.Run("test parser Next() on empty buffer", func(t *testing.T) {
		p := NewParser(f, testHeaders, "LapCurrentLapTime")
		p.current = 390

		vars, next := p.Next()
		if vars != nil {
			t.Errorf("expected vars to be nil, got %v", vars)
		}
		if next {
			t.Error("expected next to be false")
		}
	})
}

type testReader struct {
	*bytes.Reader
}

func (t testReader) Close() error                        { return nil }
func (t testReader) ReadFrom(r io.Reader) (int64, error) { return 0, nil }

func TestParserRead(t *testing.T) {
	data := make([]byte, 128)
	rand.Read(data)

	r := testReader{bytes.NewReader(data)}
	p := Parser{reader: r}

	t.Run("parser read EOF", func(t *testing.T) {
		p.length = 10
		result := p.read(0)
		expected := []byte(data[:10])
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected parsed buffer to match expected buffer.\nactual: %v\nexpected: %v", result, expected)
		}
	})

	t.Run("parser read EOF", func(t *testing.T) {

		p.length = 129
		if p.read(0) != nil {
			t.Error("expected nil")
		}
	})

}
