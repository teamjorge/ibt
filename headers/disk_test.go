package headers

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var expectedDiskHeader = DiskHeader{
	StartDate:   1719258336,
	StartTime:   932.000000635264,
	EndTime:     938.4833339685914,
	LapCount:    1,
	RecordCount: 390,
}

func TestDiskHeaders(t *testing.T) {
	t.Run("valid header file", func(t *testing.T) {
		f, err := os.Open("../.testing/valid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		// Read telemetry header to move buffer along
		if _, err := f.Seek(int64(TELEMETRY_HEADER_BYTES_SIZE), 0); err != nil {
			t.Errorf("failed to read telemetry header to move buffer along: %v", err)
		}

		output, err := ReadDiskHeader(f)
		if err != nil {
			t.Errorf("failed to parse disk header for testing file - %v", err)
			return
		}

		if !reflect.DeepEqual(*output, expectedDiskHeader) {
			t.Errorf("expected diskHeader does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedDiskHeader, *output)
		}
	})

	t.Run("invalid header file", func(t *testing.T) {
		f, err := os.Open("../.testing/invalid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		a, err := ReadDiskHeader(f)
		if err == nil {
			t.Error("expected diskHeader parsing of invalid file to return an error")
		}
		fmt.Println(a)
	})

	t.Run("empty file", func(t *testing.T) {
		f, err := os.Open("../.testing/empty_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		_, err = ReadDiskHeader(f)
		if err == nil {
			t.Error("expected diskHeader parsing of empty file to return an error")
		}
	})
}
