package ibt

import (
	"bytes"
	"crypto/rand"
	"io"
	"os"
	"reflect"
	"sort"
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

		if !reflect.DeepEqual(p.header, testHeaders) {
			t.Errorf("expected headers to be %v, received: %v", testHeaders, p.header)
		}

		expectedWhitelist := []string{"Lap", "Speed"}
		receivedWhitelist := p.whitelist
		sort.Strings(receivedWhitelist)

		if receivedWhitelist[0] != expectedWhitelist[0] || receivedWhitelist[1] != expectedWhitelist[1] {
			t.Errorf("expected whitelist to be %v, received: %v", expectedWhitelist, p.whitelist)
		}

		if len(p.header.VarHeader) != 276 {
			t.Errorf("expected varHeader to be of length %d, actual: %d", 276, len(p.header.VarHeader))
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
			if vars["LapCurrentLapTime"] != expectedValue {
				t.Errorf("expected LapCurrentLapTime value to equal %f, got %f", expectedValue, vars["LapCurrentLapTime"])
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
		if vars["LapCurrentLapTime"] != expectedValue1 {
			t.Errorf("expected LapCurrentLapTime value to equal %f, got %f", expectedValue1, vars["LapCurrentLapTime"])
		}
		if !next {
			t.Error("expected additional var values to be available after iteration")
		}

		expectedValue2 := float32(44.145233)
		vars, next = p.Next()
		if vars["LapCurrentLapTime"] != expectedValue2 {
			t.Errorf("expected LapCurrentLapTime value to equal %f, got %f", expectedValue2, vars["LapCurrentLapTime"])
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
	if _, err := rand.Read(data); err != nil {
		t.Errorf("failed to create random sequence of bytes - %v", err)
	}

	r := testReader{bytes.NewReader(data)}
	p := Parser{reader: r}

	t.Run("parser read buffer", func(t *testing.T) {
		p.header = &headers.Header{TelemetryHeader: &headers.TelemetryHeader{BufLen: 10}}
		result := p.read(0)
		expected := []byte(data[:10])
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected parsed buffer to match expected buffer.\nactual: %v\nexpected: %v", result, expected)
		}
	})

	t.Run("parser read EOF", func(t *testing.T) {
		p.header = &headers.Header{TelemetryHeader: &headers.TelemetryHeader{BufLen: 129}}
		if p.read(0) != nil {
			t.Error("expected nil")
		}
	})

}

func TestParserParseAt(t *testing.T) {
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

	t.Run("parser ParseAt buffer", func(t *testing.T) {
		p := NewParser(f, testHeaders, "LapCurrentLapTime")

		res := p.ParseAt(testHeaders.TelemetryHeader.BufOffset + (testHeaders.TelemetryHeader.BufLen * 5))

		if res["LapCurrentLapTime"] != float32(37.7452354431) {
			t.Errorf("expected parsed lap time to be %.10f. received: %.10f", 37.7452354431, res["LapCurrentLapTime"])
		}
	})

	t.Run("parser ParseAt EOF buffer", func(t *testing.T) {
		p := NewParser(f, testHeaders, "LapCurrentLapTime")

		res := p.ParseAt(testHeaders.TelemetryHeader.BufOffset + (testHeaders.TelemetryHeader.BufLen * 390))

		if res != nil {
			t.Errorf("expected parsed tick to be nil. received %v", res)
		}

	})
}

func TestSeek(t *testing.T) {
	t.Run("test seek", func(t *testing.T) {
		parser := NewParser(nil, nil)

		parser.Seek(50)

		if parser.current != 50 {
			t.Errorf("expected parser current to be %d. received %d", 50, parser.current)
		}
	})
}

func TestUpdateWhitelist(t *testing.T) {
	t.Run("parser read buffer", func(t *testing.T) {
		parser := NewParser(nil, nil, "Speed")

		parser.UpdateWhitelist("Speed", "Lap")

		if parser.whitelist[0] != "Speed" || parser.whitelist[1] != "Lap" {
			t.Errorf("expected updated white list to be %v. received: %v", []string{"Speed", "Lap"}, parser.whitelist)
		}
	})
}
