package headers

import (
	"fmt"

	"github.com/teamjorge/ibt/utilities"
)

const (
	VAR_BUFFER_HEADER_BASE_OFFSET int = 48
	VAR_BUFFER_INCREMENT          int = 16
)

// VarBuffer is a header providing information on each of the available live data buffers
type VarBuffer struct {
	TickCount int
	BufOffset int
}

// ParseVarBufferHeader retrieves the metadata of available live data buffers.
func ReadVarBufferHeaders(r Reader, numBuf int) ([]VarBuffer, error) {
	varBuffers := make([]VarBuffer, 0)
	for i := 0; i < numBuf; i++ {
		rbuf := make([]byte, 8)
		_, err := r.ReadAt(rbuf, int64(VAR_BUFFER_HEADER_BASE_OFFSET+i*VAR_BUFFER_INCREMENT))
		if err != nil {
			return nil, fmt.Errorf("failed to read VarBuffer header %d: %v", i, err)
		}
		currentVb := VarBuffer{
			utilities.Byte4ToInt(rbuf[0:4]),
			utilities.Byte4ToInt(rbuf[4:8]),
		}

		if currentVb.BufOffset == 0 || currentVb.TickCount == 0 {
			return nil, fmt.Errorf("invalid VarBuffer headers detected: %+v", currentVb)
		}

		varBuffers = append(varBuffers, currentVb)
	}

	return varBuffers, nil
}
