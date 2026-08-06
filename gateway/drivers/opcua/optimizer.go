package opcua

import "sync"

type Optimizer struct {
	lastValues sync.Map
}

func NewOptimizer() *Optimizer {
	return &Optimizer{}
}

func (o *Optimizer) LastValueFilter(sample Sample) (Sample, bool) {
	if last, ok := o.lastValues.Load(sample.ID); ok {
		if last == sample.Value {
			return sample, false
		}
	}
	o.lastValues.Store(sample.ID, sample.Value)
	return sample, true
}