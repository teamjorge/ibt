package headers

import (
	"os"
	"reflect"
	"testing"
)

var expectedHeader = Header{
	telemHeader: &expectedTelemetryHeader,
	diskHeader:  &expectedDiskHeader,
	sessionInfo: &expectedSessionInfo,
	varHeader:   expectedVarHeaders,
}

func TestHeaders(t *testing.T) {
	t.Run("valid header file", func(t *testing.T) {
		f, err := os.Open("../.testing/valid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		output, err := ParseHeaders(f)
		if err != nil {
			t.Errorf("failed to parse header for testing file - %v", err)
			return
		}

		if !reflect.DeepEqual(expectedHeader.TelemetryHeader(), output.telemHeader) {
			t.Errorf("expected header telemetryHeader does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedHeader.telemHeader, output.telemHeader)
		}

		if !reflect.DeepEqual(expectedHeader.DiskHeader(), output.diskHeader) {
			t.Errorf("expected header diskHeader does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedHeader.diskHeader, output.diskHeader)
		}

		if !reflect.DeepEqual(expectedHeader.SessionInfo(), output.sessionInfo) {
			t.Errorf("expected header sessionInfo does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedHeader.sessionInfo, output.sessionInfo)
		}

		if !reflect.DeepEqual(expectedHeader.VarHeader(), output.varHeader) {
			t.Errorf("expected header varHeader does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedHeader.varHeader, output.varHeader)
		}
	})

	t.Run("invalid header file", func(t *testing.T) {
		f, err := os.Open("../.testing/invalid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		_, err = ParseHeaders(f)
		if err == nil {
			t.Error("expected parsing of invalid file to return an error")
		}
	})

	t.Run("empty file", func(t *testing.T) {
		f, err := os.Open("../.testing/empty_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		_, err = ParseHeaders(f)
		if err == nil {
			t.Error("expected parsing of empty file to return an error")
		}
	})
}
