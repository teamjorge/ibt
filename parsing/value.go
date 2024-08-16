package parsing

import (
	"github.com/teamjorge/ibt/headers"
	"github.com/teamjorge/ibt/utilities"
)

// readVarValue extracts the telemetry variable value from the given buffer based on the provided metadata.
//
// This function will ensure that the underlying type of the value is correct.
func readVarValue(buf []byte, vh headers.VarHeader) interface{} {
	var rbuf []byte

	offset := vh.Offset
	var value interface{}

	if vh.Count > 1 {
		switch vh.Rtype {
		case 0:
			rbuf = buf[offset:vh.Count]
			res := make([]uint8, 0)
			for _, x := range rbuf[offset : offset+vh.Count] {
				res = append(res, uint8(x))
			}
			value = res
		case 1:
			rbuf = buf[offset : offset+vh.Count]
			res := make([]bool, 0)
			for _, x := range rbuf {
				res = append(res, x > 0)
			}
			value = res
		case 2:
			rbuf = buf[offset : offset+vh.Count*4]
			res := make([]int, 0)
			for i := 0; i < len(rbuf); i += 4 {
				res = append(res, utilities.Byte4ToInt(rbuf[i:i+4]))
			}
			value = res
		case 3:
			rbuf = buf[offset : offset+vh.Count*4]
			res := make([]string, 0)
			for i := 0; i < len(rbuf); i += 4 {
				res = append(res, utilities.Byte4toBitField(rbuf[i:i+4]))
			}
			value = res
		case 4:
			rbuf = buf[offset : offset+vh.Count*4]
			res := make([]float32, 0)
			for i := 0; i < len(rbuf); i += 4 {
				res = append(res, utilities.Byte4ToFloat(rbuf[i:i+4]))
			}
			value = res
		case 5:
			rbuf = buf[offset : offset+vh.Count*8]
			res := make([]float64, 0)
			for i := 0; i < len(rbuf); i += 8 {
				res = append(res, utilities.Byte8ToFloat(rbuf[i:i+8]))
			}
			value = res
		}
	} else {
		switch vh.Rtype {
		case 0:
			rbuf = buf[offset : offset+1]
			value = rbuf[0]
		case 1:
			rbuf = buf[offset : offset+1]
			value = int(rbuf[0]) > 0
		case 2:
			rbuf = buf[offset : offset+4]
			value = utilities.Byte4ToInt(rbuf)
		case 3:
			rbuf = buf[offset : offset+4]
			value = utilities.Byte4toBitField(rbuf)
		case 4:
			rbuf = buf[offset : offset+4]
			value = utilities.Byte4ToFloat(rbuf)
		case 5:
			rbuf = buf[offset : offset+8]
			value = utilities.Byte8ToFloat(rbuf)
		}
	}

	return value
}
