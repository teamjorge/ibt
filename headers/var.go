package headers

import (
	"fmt"
	"unicode/utf8"

	"github.com/teamjorge/ibt/utilities"
	"golang.org/x/exp/maps"
)

const (
	VAR_HEADER_BYTES_SIZE int = 144
)

// VarHeader provides the available telemetry variable information and values.
type VarHeader struct {
	// Rtype is the variable value type.
	//
	// Possible values:
	// 0: String
	// 1: Boolean
	// 2: Int
	// 3: Byte
	// 4: Float32
	// 5: Float64
	Rtype int `json:"rtype,omitempty"`
	// Offset in the buffer where the value can be found
	Offset int `json:"offset,omitempty"`
	// Number of items the value consists of. >1 means it is an array.
	Count int `json:"count,omitempty"`
	// Indicates if the value is time based
	CountAsTime bool `json:"count_as_time,omitempty"`
	// Description of the variable
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	// Unit of measurement for the variable value
	Unit string `json:"unit,omitempty"`

	// Value of the variable. The is parsed during iteration of telemetry data.
	Value interface{} `json:"value"`
}

// ReadVarHeader populates the the VarHeader with the necessary metadata.
//
// This function will not populate the value field, but rather provides a template for retrieving values during
// telemetry processing.
//
// Validation is performed by ensuring all variable names conform the UTF-8.
func ReadVarHeader(reader Reader, numVars, offset int) (map[string]VarHeader, error) {
	varHeaderBuf := make([]byte, numVars*VAR_HEADER_BYTES_SIZE)

	_, err := reader.ReadAt(varHeaderBuf, int64(offset))
	if err != nil {
		return nil, err
	}

	varHeaders := make(map[string]VarHeader, 0)

	start := 0
	for i := 0; i < numVars; i++ {
		h := VarHeader{
			Rtype:       utilities.Byte4ToInt(varHeaderBuf[start+0 : start+4]),
			Offset:      utilities.Byte4ToInt(varHeaderBuf[start+4 : start+8]),
			Count:       utilities.Byte4ToInt(varHeaderBuf[start+8 : start+12]),
			CountAsTime: int(varHeaderBuf[start+12]) > 0,
			// Padded
			Name:        utilities.BytesToString(varHeaderBuf[start+16 : start+48]),
			Description: utilities.BytesToString(varHeaderBuf[start+48 : start+112]),
			Unit:        utilities.BytesToString(varHeaderBuf[start+112 : start+144]),
		}
		start += VAR_HEADER_BYTES_SIZE

		varHeaders[h.Name] = h

		if !utf8.Valid([]byte(h.Name)) {
			return nil, fmt.Errorf("invalid vars detected at item %d. current var name: %s", i, h.Name)
		}
	}

	return varHeaders, nil
}

// AvailableVars for each tick of telemetry data.
//
// This is useful when determining which variables are available for a specific car.
func AvailableVars(varHeaders map[string]VarHeader) []string {
	return maps.Keys(varHeaders)
}
