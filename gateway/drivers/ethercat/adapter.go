package ethercat

import "context"

// Adapter is the interface that abstracts the underlying EtherCAT master implementation.
// This allows for swapping the low-level driver (e.g., SOEM, native Go) without
// changing the high-level business logic.
type Adapter interface {
	// Connect initializes the EtherCAT master on the specified network interface.
	Connect(ctx context.Context, cfg ConnectionConfig) error

	// Disconnect closes the master and releases the network interface.
	Disconnect(ctx context.Context) error

	// IsConnected returns true if the adapter is connected and the master is running.
	IsConnected() bool

	// ScanSlaves scans the bus and returns information about all discovered slaves.
	ScanSlaves(ctx context.Context) ([]SlaveInfo, error)

	// ConfigurePDOs sets up the PDO mapping for all slaves. This must be called
	// before the cyclic exchange can start.
	ConfigurePDOs(ctx context.Context, mappings []PDOMapping) error

	// ReadPDOs executes one cycle of reading all input PDOs from the slaves.
	// It returns a map where the key is the slave index and the value is the raw PDO data.
	ReadPDOs(ctx context.Context) (map[uint16][]byte, error)

	// WritePDOs executes one cycle of writing all output PDOs to the slaves.
	// The input is a map where the key is the slave index and the value is the raw PDO data.
	WritePDOs(ctx context.Context, data map[uint16][]byte) error

	// ReadSDO performs a service data object read from a specific slave.
	ReadSDO(ctx context.Context, req SDORequest) (SDOResponse, error)

	// WriteSDO performs a service data object write to a specific slave.
	WriteSDO(ctx context.Context, req SDORequest) error
}