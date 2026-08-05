// drivers/bacnet/optimizer.go
package bacnet

import (
	"fmt"
	"sync"
)

// Optimizer processes raw BACnet samples to reduce redundancy or apply business logic.
// It can filter out unchanged values, aggregate data, or perform other optimizations
// before samples are passed to higher layers.
type Optimizer struct {
	// For now, a simple last-value cache to filter out unchanged samples.
	lastValues map[string]interface{}
	mu         sync.RWMutex
}

// NewOptimizer creates a new Optimizer instance.
func NewOptimizer() *Optimizer {
	return &Optimizer{
		lastValues: make(map[string]interface{}),
	}
}

// ProcessSamples takes a slice of raw samples and returns an optimized slice.
// This is a basic implementation that only passes through samples whose value has changed.
// More sophisticated logic (e.g., deadband, aggregation) can be added here.
func (o *Optimizer) ProcessSamples(rawSamples []Sample) []Sample {
	o.mu.Lock()
	defer o.mu.Unlock()

	optimizedSamples := make([]Sample, 0, len(rawSamples))

	for _, sample := range rawSamples {
		lastValue, exists := o.lastValues[sample.ID]

		// Only add the sample if its value has changed or if it's a new signal.
		// Also, always pass through bad quality samples or errors.
		if !exists || !valuesEqual(lastValue, sample.Value) || sample.Quality != QualityGood {
			optimizedSamples = append(optimizedSamples, sample)
			o.lastValues[sample.ID] = sample.Value // Update last known good value
		}
	}

	return optimizedSamples
}

// valuesEqual performs a deep comparison for interface{} values.
// This is a simplified comparison; for complex types, a more robust deep equality check
// or custom comparison logic might be needed.
func valuesEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	// Basic comparison for primitive types
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
