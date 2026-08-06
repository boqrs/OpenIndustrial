package ethernetip

import (
	"context"
	"log"
	"sync"
)

// Optimizer filters out duplicate samples, only passing on changes.
type Optimizer struct {
	input     <-chan Sample
	output    chan<- Sample
	lastState map[string]interface{}
	mu        sync.RWMutex
}

// NewOptimizer creates a new optimizer.
func NewOptimizer(input <-chan Sample, output chan<- Sample) *Optimizer {
	return &Optimizer{
		input:     input,
		output:    output,
		lastState: make(map[string]interface{}),
	}
}

// Start begins the optimization loop, processing samples from the input channel.
func (o *Optimizer) Start(ctx context.Context) {
	log.Println("Starting EtherNet/IP optimizer...")
	for {
		select {
		case sample := <-o.input:
			o.process(sample)
		case <-ctx.Done():
			log.Println("Stopping EtherNet/IP optimizer.")
			return
		}
	}
}

// process checks if a sample's value has changed and, if so, forwards it.
func (o *Optimizer) process(sample Sample) {
	o.mu.RLock()
	lastValue, found := o.lastState[sample.PointID]
	o.mu.RUnlock()

	// Always forward bad quality samples to indicate problems.
	if sample.Quality != QualityGood {
		o.updateAndSend(sample)
		return
	}

	// If the value has changed or it's the first time we see this point,
	// update the state and forward the sample.
	if !found || lastValue != sample.Value {
		o.updateAndSend(sample)
	}
}

func (o *Optimizer) updateAndSend(sample Sample) {
	o.mu.Lock()
	o.lastState[sample.PointID] = sample.Value
	o.mu.Unlock()
	o.output <- sample
}