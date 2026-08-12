package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Envelope 是我们系统中所有事件的标准化结构。
// 它包含了路由、追踪、和解析事件所需的所有元数据。
type Envelope struct {
	ID             string          `json:"id"`
	Type           string          `json:"type"`
	Version        int             `json:"version"`
	TenantID       string          `json:"tenant_id"`
	AggregateType  string          `json:"aggregate_type"`
	AggregateID    string          `json:"aggregate_id"`
	OccurredAt     time.Time       `json:"occurred_at"`
	TraceID        string          `json:"trace_id,omitempty"`
	Payload        json.RawMessage `json:"payload"`
}

// NewEnvelope 是一个工厂函数，用于方便地创建新的事件信封。
func NewEnvelope(eventType, aggregateType, aggregateID, tenantID string, payload interface{}) (*Envelope, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Envelope{
		ID:            "evt_" + uuid.New().String(),
		Type:          eventType,
		Version:       1,
		TenantID:      tenantID,
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		OccurredAt:    time.Now().UTC(),
		Payload:       payloadBytes,
	}, nil
}