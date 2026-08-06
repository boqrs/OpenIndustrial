package can

// Config holds the complete configuration for the CAN driver.
type Config struct {
	// Name is a user-defined name for this driver instance.
	Name string `json:"name"`
	// Connection settings for the CAN adapter.
	Connection ConnectionConfig `json:"connection"`
	// Bitrate of the CAN bus. Note: This is often configured at the OS level.
	Bitrate int `json:"bitrate"`
	// FD indicates if CAN FD (Flexible Data-rate) is enabled.
	FD bool `json:"fd"`
	// Signals is a list of all signals to be decoded from incoming frames.
	Signals []Signal `json:"signals"`
	// TxMessages is a list of messages to be transmitted periodically by the poller.
	TxMessages []TxConfig `json:"txMessages"`
}