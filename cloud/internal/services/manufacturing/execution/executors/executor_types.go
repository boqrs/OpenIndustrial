package executors

const (
	OperationTypeSNGenerate       = "SN_GENERATE"
	OperationTypeCertificateIssue = "CERTIFICATE_ISSUE"
	OperationTypeSNBind           = "SN_BIND"
	OperationTypeCertificateBind  = "CERTIFICATE_BIND"
	OperationTypeDeviceBind       = "DEVICE_BIND"
	OperationTypeQualityCheck     = "QUALITY_CHECK"
	OperationTypeDataCollection   = "DATA_COLLECTION"
)

// OperationInput defines the data provided to an executor.
type OperationInput struct {
	ExecutionID          uint
	ExecutionOperationID uint
	WorkOrderID          uint
	ProductID            uint
	DeviceID             *uint
	Parameters           map[string]any
}

// OperationOutput defines the data returned by an executor.
type OperationOutput struct {
	// Result is the structured data produced by the operation,
	// which will be persisted for traceability.
	Result map[string]any
	// References are pointers to other system resources created or
	// affected by the operation, but are not part of the direct result.
	References map[string]any
}
