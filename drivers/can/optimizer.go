package can

import "sync"

// Optimizer filters out redundant (unchanged) signal values.
type Optimizer struct {
	lastValues map[string]any
	mu         sync.Mutex
}

// NewOptimizer creates a new optimizer.
func NewOptimizer() *Optimizer {
	return &Optimizer{
		lastValues: make(map[string]any),
	}
}

// Optimize checks if a new signal value is different from the last recorded one.
// It returns true if the value is new or has changed, otherwise false.
func (o *Optimizer) Optimize(value SignalValue) bool {
	o.mu.Lock()
	defer o.mu.Unlock()

	last, ok := o.lastValues[value.ID]
	if !ok || last != value.Value {
		o.lastValues[value.ID] = value.Value
		return true
	}

	return false
}