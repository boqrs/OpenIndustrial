package ethercat

// This file is intended to hold helper functions for SDO requests and responses,
// for example, functions to encode/decode specific data types for SDO transfers.

// For now, most of the SDO logic is defined in the Adapter interface and will be
// implemented in adapter_soem.go, as it involves direct CGO calls.

// Example helper function (can be expanded later):

// NewReadRequest creates a simple SDO read request.
func NewReadRequest(slave uint16, index uint16, subIndex uint8) SDORequest {
	return SDORequest{
		Slave:    slave,
		Index:    index,
		SubIndex: subIndex,
	}
}

// NewWriteRequest creates a simple SDO write request.
func NewWriteRequest(slave uint16, index uint16, subIndex uint8, value interface{}) SDORequest {
	return SDORequest{
		Slave:    slave,
		Index:    index,
		SubIndex: subIndex,
		Value:    value,
	}
}