package executors

const (
	// SN_GENERATE is kept for cloud-generated serial number scenarios.
	OperationTypeSNGenerate = "SN_GENERATE"

	// SN_WRITE records a serial number provided by an external system.
	OperationTypeSNWrite = "SN_WRITE"

	// CA_ISSUE issues a device certificate through the configured CA.
	OperationTypeCAIssue = "CA_ISSUE"

	// FUNCTION_TEST represents a production functional test.
	OperationTypeFunctionTest = "FUNCTION_TEST"

	// Legacy / future operation types.
	OperationTypeSNBind          = "SN_BIND"
	OperationTypeCertificateBind = "CERTIFICATE_BIND"
	OperationTypeDeviceBind      = "DEVICE_BIND"
	OperationTypeQualityCheck    = "QUALITY_CHECK"
	OperationTypeDataCollection  = "DATA_COLLECTION"
)

// OperationInput defines the data provided to an executor.
type OperationInput struct {
	ExecutionID          uint
	ExecutionOperationID uint
	WorkOrderID          uint
	ProductID            uint
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
