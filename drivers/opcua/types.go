package opcua

import "time"

type Quality string

const (
	QualityGood    Quality = "Good"
	QualityBad     Quality = "Bad"
	QualityUncertain Quality = "Uncertain"
)

type NodeMapping struct {
	ID     string `json:"id"`
	NodeID string `json:"nodeId"`
}

type Sample struct {
	ID        string      `json:"id"`
	NodeID    string      `json:"nodeId"`
	Value     interface{} `json:"value"`
	Timestamp time.Time   `json:"timestamp"`
	Quality   Quality     `json:"quality"`
}