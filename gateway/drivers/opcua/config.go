package opcua

import (
	"fmt"
	"time"
)

type ConnectionConfig struct {
	EndpointURL      string        `json:"endpointUrl"`
	SecurityPolicy   string        `json:"securityPolicy"`
	SecurityMode     string        `json:"securityMode"`
	AuthType         string        `json:"authType"`
	Username         string        `json:"username"`
	Password         string        `json:"password"`
	Timeout          time.Duration `json:"timeout"`
	CertPath         string        `json:"certPath"`
	KeyPath          string        `json:"keyPath"`
}

func (c *ConnectionConfig) Validate() error {
	if c.EndpointURL == "" {
		return fmt.Errorf("endpointUrl is required")
	}
	return nil
}

type SubscriptionConfig struct {
	PublishingInterval time.Duration `json:"publishingInterval"`
	LifetimeCount      uint32        `json:"lifetimeCount"`
	MaxKeepAliveCount  uint32        `json:"maxKeepAliveCount"`
	SamplingInterval   time.Duration `json:"samplingInterval"`
}


type PollConfig struct {
	Interval time.Duration `json:"interval"`
}