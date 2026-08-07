package data

// RuleAction defines the action to be taken when a stream rule's condition is met.
type RuleAction string

const (
	ActionAlarm      RuleAction = "alarm"
	ActionEvent      RuleAction = "event"
	ActionCommand    RuleAction = "command"
	ActionNotification RuleAction = "notification"
)

// StreamRule defines a condition to be evaluated against the real-time data stream.
// e.g., "If motor_current > 5A for 10 seconds, then trigger a CRITICAL alarm".
type StreamRule struct {
	ID         string     `json:"id"`
	OrgID      string     `json:"org_id"`
	Metric     string     `json:"metric"`
	// Expression is a logical condition, e.g., "> 5", "== 'ERROR'".
	Expression string     `json:"expression"`
	// Window defines a time duration for the rule, e.g., "10s", "5m".
	Window     string     `json:"window,omitempty"`
	Action     RuleAction `json:"action"`
	// ActionPayload provides context for the action, e.g., alarm message or command parameters.
	ActionPayload map[string]interface{} `json:"action_payload"`
	IsEnabled  bool       `json:"is_enabled"`
}