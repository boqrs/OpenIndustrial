package codec

import (
	"fmt"
	"sync"
)

// DataType represents the type of data to be decoded.
type DataType string

const (
	DataTypeBool    DataType = "bool"
	DataTypeInt16   DataType = "int16"
	DataTypeUint16  DataType = "uint16"
	DataTypeInt32   DataType = "int32"
	DataTypeUint32  DataType = "uint32"
	DataTypeFloat32 DataType = "float32"
	DataTypeFloat64 DataType = "float64"
	DataTypeInt64   DataType = "int64"
	DataTypeUint64  DataType = "uint64"
	DataTypeBytes   DataType = "bytes"
	DataTypeString  DataType = "string"
	DataTypeBCD     DataType = "bcd"
)

// DecoderFunc is a function that decodes raw bytes into a specific Go type.
// It takes the raw byte slice, byte order, and word order as input.
type DecoderFunc func(data []byte, byteOrder Endian, wordOrder WordOrder) (interface{}, error)

var (
	decoderRegistry = make(map[DataType]DecoderFunc)
	registryMutex   sync.RWMutex
)

// RegisterDecoder registers a DecoderFunc for a given DataType.
// It is safe for concurrent use.
func RegisterDecoder(dataType DataType, decoderFn DecoderFunc) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	if _, exists := decoderRegistry[dataType]; exists {
		// Log or handle error if a decoder for this type is already registered
		fmt.Printf("Warning: Decoder for DataType %s already registered. Overwriting.\n", dataType)
	}
	decoderRegistry[dataType] = decoderFn
}

// GetDecoder retrieves the DecoderFunc for a given DataType.
// It is safe for concurrent use.
func GetDecoder(dataType DataType) (DecoderFunc, error) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	decoderFn, ok := decoderRegistry[dataType]
	if !ok {
		return nil, fmt.Errorf("no decoder registered for DataType: %s", dataType)
	}
	return decoderFn, nil
}