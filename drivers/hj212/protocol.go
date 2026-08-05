package hj212

import (
	"strings"
)

// parseDataSegment takes the raw string from the "CP" field of a message
// and parses it into a structured DataSegment object.
// Example input: "QN=...;ST=...;CN=...;PW=...;MN=...;w01018-Rtd=56.3,w21003-Rtd=7.2"
func ParseDataSegment(cpData string) (*DataSegment, error) {
	segment := &DataSegment{
		Pollutants: make(map[string]string),
	}

	pairs := strings.Split(cpData, ";")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue // Ignore malformed pairs
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "QN":
			segment.QN = value
		case "ST":
			segment.ST = value
		case "CN":
			segment.CN = value
		case "PW":
			segment.PW = value
		case "MN":
			segment.MN = value
		default:
			// All other keys are assumed to be pollutant data
			segment.Pollutants[key] = value
		}
	}

	// Basic validation
	if segment.QN == "" || segment.MN == "" || segment.CN == "" {
		return nil, ErrMissingRequiredField
	}

	return segment, nil
}

// buildDataSegment serializes a DataSegment object back into a "CP" string.
func BuildDataSegment(segment *DataSegment) string {
	var parts []string
	parts = append(parts, "QN="+segment.QN)
	parts = append(parts, "ST="+segment.ST)
	parts = append(parts, "CN="+segment.CN)
	parts = append(parts, "PW="+segment.PW)
	parts = append(parts, "MN="+segment.MN)

	for key, value := range segment.Pollutants {
		parts = append(parts, key+"="+value)
	}

	return strings.Join(parts, ";")
}