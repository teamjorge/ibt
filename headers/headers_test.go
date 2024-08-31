package headers

import (
	"bytes"
	"fmt"
)

// Shared resources for headers tests

type mockReader struct {
	*bytes.Reader
}

// Creates a new mockReader
//
// mockReader is used to scramble the contents of a buffer.
//
// Scrambling is essentially replacing a subset of bytes with 0x0. This is useful when
// testing for specific errors that might occur without having to create specific files for it.
//
//   - bufSize - Size of the buffer you want the reader to have
//   - start - Start offset for scrambling
//   - end - End offset for scrambling
func newMockReader(reader Reader, bufSize, start, end int) (*mockReader, error) {
	buf := make([]byte, bufSize)

	_, err := reader.ReadAt(buf, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to read into mockReader buf: %v", err)
	}

	newBuf := buf[:start]
	newBuf = append(newBuf, make([]byte, end-start)...)
	newBuf = append(newBuf, buf[end:]...)

	return &mockReader{bytes.NewReader(newBuf)}, nil
}

func (m mockReader) Close() error { return nil }
