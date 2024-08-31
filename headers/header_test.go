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

	t.Run("test fail telem header", func(t *testing.T) {
		f, err := os.Open("../.testing/valid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		mock, err := newMockReader(
			f,
			TELEMETRY_HEADER_BYTES_SIZE,
			0,
			8,
		)
		if err != nil {
			t.Errorf("failed to initialize a test mockReader - %v", err)
		}

		h, err := ParseHeaders(mock)
		if err != nil && err.Error()[:33] != "failed to parse telemetry header:" {
			t.Errorf("expected failed ParseHeaders err to be %v. received: %v. full error: %v",
				"failed to parse telemetry header:", err.Error()[:33], err)
		}

		if err == nil {
			t.Errorf("expected an error to occur when parsing a defective telemetry header. received telemetry header: %v", h.TelemetryHeader)
		}
	})

	t.Run("test fail disk header", func(t *testing.T) {
		f, err := os.Open("../.testing/valid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		mock, err := newMockReader(
			f,
			TELEMETRY_HEADER_BYTES_SIZE+DISK_HEADER_BYTES_SIZE,
			TELEMETRY_HEADER_BYTES_SIZE,
			TELEMETRY_HEADER_BYTES_SIZE+8,
		)
		if err != nil {
			t.Errorf("failed to initialize a test mockReader - %v", err)
		}

		h, err := ParseHeaders(mock)
		if err != nil && err.Error()[:28] != "failed to parse disk header:" {
			t.Errorf("expected failed ParseHeaders err to be %v. received: %v. full error: %v",
				"failed to parse disk header:", err.Error()[:28], err)
		}

		if err == nil {
			t.Errorf("expected an error to occur when parsing a defective disk header. received disk header: %v", h.DiskHeader)
		}
	})

	t.Run("test fail varheader", func(t *testing.T) {
		f, err := os.Open("../.testing/valid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		mock, err := newMockReader(
			f,
			expectedTelemetryHeader.VarHeaderOffset+VAR_HEADER_BYTES_SIZE,
			expectedTelemetryHeader.VarHeaderOffset,
			expectedTelemetryHeader.VarHeaderOffset+VAR_HEADER_BYTES_SIZE,
		)
		if err != nil {
			t.Errorf("failed to initialize a test mockReader - %v", err)
		}

		h, err := ParseHeaders(mock)
		if err != nil && err.Error()[:32] != "failed to parse variable header:" {
			t.Errorf("expected failed ParseHeaders err to be %v. received: %v. full error: %v",
				"failed to parse variable header:", err.Error()[:32], err)
		}

		if err == nil {
			t.Errorf("expected an error to occur when parsing a defective varheader. received varheader: %v", h.VarHeader)
		}
	})

	t.Run("test fail varbuffer", func(t *testing.T) {
		f, err := os.Open("../.testing/valid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		mock, err := newMockReader(
			f,
			expectedTelemetryHeader.VarHeaderOffset+(VAR_HEADER_BYTES_SIZE*expectedTelemetryHeader.NumVars),
			VAR_BUFFER_HEADER_BASE_OFFSET,
			VAR_BUFFER_HEADER_BASE_OFFSET+8,
		)
		if err != nil {
			t.Errorf("failed to initialize a test mockReader - %v", err)
		}

		h, err := ParseHeaders(mock)
		if err != nil && err.Error()[:34] != "failed to parse var buffer header:" {
			t.Errorf("expected failed ParseHeaders err to be %v. received: %v. full error: %v",
				"failed to parse var buffer header:", err.Error()[:34], err)
		}

		if err == nil {
			t.Errorf("expected an error to occur when parsing a defective varheader. received varheader: %v", h.VarHeader)
		}
	})

	t.Run("test fail session info", func(t *testing.T) {
		f, err := os.Open("../.testing/valid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		mock, err := newMockReader(
			f,
			expectedTelemetryHeader.SessionInfoOffset+expectedTelemetryHeader.SessionInfoLength,
			expectedTelemetryHeader.SessionInfoOffset,
			expectedTelemetryHeader.SessionInfoOffset+expectedTelemetryHeader.SessionInfoLength,
		)
		if err != nil {
			t.Errorf("failed to initialize a test mockReader - %v", err)
		}

		h, err := ParseHeaders(mock)
		if err != nil && err.Error()[:29] != "failed to parse session info:" {
			t.Errorf("expected failed ParseHeaders err to be %v. received: %v. full error: %v",
				"failed to parse session info:", err.Error()[:29], err)
		}

		if err == nil {
			t.Errorf("expected an error to occur when parsing a defective session info. received session info: %v", h.DiskHeader)
		}
	})
}

func TestUpdateVarBuffer(t *testing.T) {
	t.Run("test valid update var buffer", func(t *testing.T) {
		validIbtFile, err := os.Open("../.testing/valid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer validIbtFile.Close()

		output, err := ParseHeaders(validIbtFile)
		if err != nil {
			t.Errorf("failed to parse header for testing file - %v", err)
			return
		}

		if !reflect.DeepEqual(output.VarBuffers[0], expectedVarBuffers[0]) {
			t.Errorf("expected initial VarBuffer to be %v. received: %v", output.VarBuffers[0], expectedVarBuffers[0])
		}

		liveIbtFile, err := os.Open("../.testing/live_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer liveIbtFile.Close()

		output.TelemetryHeader.NumBuf = 3

		if err := output.UpdateVarBuffer(liveIbtFile); err != nil {
			t.Errorf("did not expect an error when updating VarBuffer. received: %v", err)
		}

		if !reflect.DeepEqual(output.VarBuffers[0], expectedLiveVarBuffers[0]) || !reflect.DeepEqual(output.VarBuffers[1], expectedLiveVarBuffers[1]) ||
			!reflect.DeepEqual(output.VarBuffers[2], expectedLiveVarBuffers[2]) {
			t.Errorf("expected updated VarBuffers to be %v. received: %v", output.VarBuffers, expectedLiveVarBuffers)
		}
	})

	t.Run("invalid header file", func(t *testing.T) {
		f, err := os.Open("../.testing/invalid_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		h := &Header{TelemetryHeader: &expectedTelemetryHeader}
		if err := h.UpdateVarBuffer(f); err == nil {
			t.Error("expected telemHeader parsing of invalid file to return an error")
		}

		if h.VarBuffers != nil {
			t.Errorf("expected var buffers to be nil. received: %v", h.VarBuffers)
		}
	})

	t.Run("empty file", func(t *testing.T) {
		f, err := os.Open("../.testing/empty_test_file.ibt")
		if err != nil {
			t.Errorf("failed to open testing file - %v", err)
			return
		}
		defer f.Close()

		h := &Header{TelemetryHeader: &expectedTelemetryHeader}
		if err := h.UpdateVarBuffer(f); err == nil {
			t.Error("expected telemHeader parsing of invalid file to return an error")
		}

		if h.VarBuffers != nil {
			t.Errorf("expected var buffers to be nil. received: %v", h.VarBuffers)
		}
	})
}
