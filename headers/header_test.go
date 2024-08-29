package headers

import (
	"os"
	"reflect"
	"testing"
)

var expectedHeader = Header{
	TelemetryHeader: &expectedTelemetryHeader,
	DiskHeader:      &expectedDiskHeader,
	VarHeader:       expectedVarHeaders,
	SessionInfo:     &expectedSessionInfo,
	VarBuffers:      expectedVarBuffers,
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

		if !reflect.DeepEqual(expectedHeader.TelemetryHeader, output.TelemetryHeader) {
			t.Errorf("expected header telemetryHeader does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedHeader.TelemetryHeader, output.TelemetryHeader)
		}

		if !reflect.DeepEqual(expectedHeader.DiskHeader, output.DiskHeader) {
			t.Errorf("expected header diskHeader does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedHeader.DiskHeader, output.DiskHeader)
		}

		if !reflect.DeepEqual(expectedHeader.SessionInfo, output.SessionInfo) {
			t.Errorf("expected header sessionInfo does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedHeader.SessionInfo, output.SessionInfo)
		}

		if !reflect.DeepEqual(expectedHeader.VarHeader, output.VarHeader) {
			t.Errorf("expected header varHeader does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedHeader.VarHeader, output.VarHeader)
		}

		if !reflect.DeepEqual(expectedHeader.VarBuffers[0], output.VarBuffers[0]) {
			t.Errorf("expected varBuffer header does not match actual. \nexpected: %+v\n \nactual: %+v\n", expectedHeader.VarBuffers[0], output.VarBuffers[0])
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
