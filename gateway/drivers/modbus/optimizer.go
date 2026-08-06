package modbus

import (
	"sort"
)

// Optimizer is responsible for optimizing Modbus read operations by grouping
// multiple PointMappings into fewer, larger Modbus requests (Batches).
type Optimizer struct{}

// NewOptimizer creates a new Optimizer instance.
func NewOptimizer() *Optimizer {
	return &Optimizer{}
}

// Optimize takes a slice of PointMapping and returns a slice of Batch.
// It groups PointMappings by RegisterType and then by contiguous address ranges.
func (o *Optimizer) Optimize(mappings []NodeMapping) []Batch {
	if len(mappings) == 0 {
		return nil
	}

	// Group mappings by RegisterType
	groupedByRegister := make(map[RegisterType][]NodeMapping)
	for _, m := range mappings {
		groupedByRegister[m.Register] = append(groupedByRegister[m.Register], m)
	}

	var batches []Batch

	// For each RegisterType, sort and then group by contiguous addresses
	for regType, regMappings := range groupedByRegister {
		// Sort mappings by address to easily find contiguous blocks
		sort.Slice(regMappings, func(i, j int) bool {
			return regMappings[i].Address < regMappings[j].Address
		})

		if len(regMappings) == 0 {
			continue
		}

		currentBatch := Batch{
			Register: regType,
			Start:    regMappings[0].Address,
			Count:    regMappings[0].Length,
			Points:   []NodeMapping{regMappings[0]},
		}

		for i := 1; i < len(regMappings); i++ {
			m := regMappings[i]
			// Check if the current mapping is contiguous with the current batch
			// A mapping is contiguous if its address is immediately after the current batch's end address
			// and its length is also considered.
			// For simplicity, we'll assume a simple contiguous check for now.
			// A more robust optimizer might consider maximum batch size (e.g., 125 registers)
			// and also the actual data length of each point.
			expectedNextAddress := currentBatch.Start + currentBatch.Count
			if m.Address == expectedNextAddress {
				// If contiguous, extend the current batch
				currentBatch.Count += m.Length
				currentBatch.Points = append(currentBatch.Points, m)
			} else {
				// Not contiguous, start a new batch
				batches = append(batches, currentBatch)
				currentBatch = Batch{
					Register: regType,
					Start:    m.Address,
					Count:    m.Length,
					Points:   []NodeMapping{m},
				}
			}
		}
		batches = append(batches, currentBatch) // Add the last batch
	}

	return batches
}