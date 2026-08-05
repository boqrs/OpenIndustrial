package ethercat

import (
	"encoding/binary"
	"log"
	"math"
)

// pdoDecoder is responsible for decoding raw PDO byte slices into structured Samples.
type pdoDecoder struct {
	inputMappings map[uint16][]PDOEntry
}

// newPdoDecoder creates a decoder based on the provided PDO mappings.
func newPdoDecoder(mappings []PDOMapping) *pdoDecoder {
	inputMaps := make(map[uint16][]PDOEntry)
	for _, m := range mappings {
		if m.Direction == PDOInput {
			inputMaps[m.Slave] = m.Entries
		}
	}
	return &pdoDecoder{inputMappings: inputMaps}
}

// decode processes the raw PDO data from all slaves and converts it into samples.
func (d *pdoDecoder) decode(rawData map[uint16][]byte) []Sample {
	var samples []Sample
	for slaveIdx, dataBytes := range rawData {
		if entries, ok := d.inputMappings[slaveIdx]; ok {
			for _, entry := range entries {
				value, err := d.extractValue(dataBytes, entry)
				if err != nil {
					log.Printf("Failed to decode PDO for slave %d, entry %s: %v", slaveIdx, entry.Name, err)
					continue
				}
				samples = append(samples, Sample{
					PointID: entry.Name, // Use the entry name as the point ID
					Value:   value,
				})
			}
		}
	}
	return samples
}

// extractValue pulls a single value from the byte slice based on the PDOEntry definition.
func (d *pdoDecoder) extractValue(data []byte, entry PDOEntry) (interface{}, error) {
	byteOffset := entry.Offset / 8
	bitOffset := entry.Offset % 8

	// Ensure we have enough bytes
	requiredBytes := (entry.Offset + entry.BitLength + 7) / 8
	if uint(len(data)) < requiredBytes {
		return nil, log.Output(2, "not enough data for PDO entry")
	}

	// This is a simplified decoder. A real implementation needs to handle bit-level access
	// and various data types more robustly.
	switch entry.DataType {
	case TypeBOOL:
		byteVal := data[byteOffset]
		return (byteVal & (1 << bitOffset)) != 0, nil
	case TypeSINT:
		return int8(data[byteOffset]), nil
	case TypeUSINT:
		return data[byteOffset], nil
	case TypeINT:
		return int16(binary.LittleEndian.Uint16(data[byteOffset:])), nil
	case TypeUINT:
		return binary.LittleEndian.Uint16(data[byteOffset:]), nil
	case TypeDINT:
		return int32(binary.LittleEndian.Uint32(data[byteOffset:])), nil
	case TypeUDINT:
		return binary.LittleEndian.Uint32(data[byteOffset:]), nil
	case TypeREAL:
		bits := binary.LittleEndian.Uint32(data[byteOffset:])
		return math.Float32frombits(bits), nil
	case TypeLREAL:
		bits := binary.LittleEndian.Uint64(data[byteOffset:])
		return math.Float64frombits(bits), nil
	default:
		return nil, log.Output(2, "unsupported data type in PDO decoder")
	}
}