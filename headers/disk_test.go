package headers

import (
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
	validF, err := os.Open("../.testing/valid_test_file.ibt")
	if err != nil {
		t.Errorf("failed to open testing file - %v", err)
		return
	}
	defer validF.Close()

	invalidF, err := os.Open("../.testing/invalid_test_file.ibt")
	if err != nil {
		t.Errorf("failed to open testing file - %v", err)
		return
	}
	defer invalidF.Close()

	t.Run("valid header file", func(t *testing.T) {
		// Read telemetry header to move buffer along
		if _, err := validF.Seek(int64(TELEMETRY_HEADER_BYTES_SIZE), 0); err != nil {
			t.Errorf("failed to read telemetry header to move buffer along: %v", err)
		}

		output, err := ReadDiskHeader(validF)
		if err != nil {
			t.Errorf("failed to parse disk header for testing file - %v", err)
			return
		}

		if !reflect.DeepEqual(*output, expectedDiskHeader) {
			t.Errorf("expected diskHeader does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedDiskHeader, *output)
		}
	})

	t.Run("invalid header file invalid end/start time", func(t *testing.T) {
		_, err := ReadDiskHeader(invalidF)

		if err != nil && err.Error()[:46] != "invalid disk header detected. values received:" {
			t.Errorf("expected invalid disk header detected error. received: %v", err.Error()[:46])
		}

		if err == nil {
			t.Error("expected diskHeader parsing of invalid file to return an error")
		}
	})

	t.Run("invalid header file invalid parsed date", func(t *testing.T) {
		mock, err := newMockReader(
			validF,
			TELEMETRY_HEADER_BYTES_SIZE+DISK_HEADER_BYTES_SIZE,
			TELEMETRY_HEADER_BYTES_SIZE,
			TELEMETRY_HEADER_BYTES_SIZE+8,
		)
		if err != nil {
			t.Errorf("failed to initialize a test mockReader - %v", err)
		}

		_, err = ReadDiskHeader(mock)

		if err != nil && err.Error()[:27] != "invalid StartDate detected:" {
			t.Errorf("expected invalid StartDate error. received: %v. full error: %v", err.Error()[:27], err)
		}

		if err == nil {
			t.Error("expected diskHeader parsing of invalid file to return an error")
		}
	})

	t.Run("empty file", func(t *testing.T) {
		f, err := os.Open("../.testing/empty_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		_, err = ReadDiskHeader(f)

		if err != nil && err.Error()[:34] != "failed to read disk header buffer:" {
			t.Errorf("expected read disk header error. received: %v. full error: %v", err.Error()[:34], err)
		}

		if err == nil {
			t.Error("expected diskHeader parsing of empty file to return an error")
		}
	})
}
