package iec104

import "sync"

// Cache is a thread-safe in-memory store for IEC 104 data points (Samples).
type Cache struct {
	mu      sync.RWMutex
	samples map[string]Sample
}

// NewCache creates and returns a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		samples: make(map[string]Sample),
	}
}

// Set adds or updates a Sample in the cache under a given point ID.
// It is safe for concurrent use.
func (c *Cache) Set(pointID string, sample Sample) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.samples[pointID] = sample
}

// Get retrieves a Sample from the cache by its point ID.
// It returns the sample and a boolean indicating whether the point ID was found.
// It is safe for concurrent use.
func (c *Cache) Get(pointID string) (Sample, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	sample, found := c.samples[pointID]
	return sample, found
}

// GetAll returns a copy of all samples currently in the cache.
// This is useful for exposing the current state without allowing direct
// modification of the internal cache map.
func (c *Cache) GetAll() map[string]Sample {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Create a copy to avoid race conditions on the returned map
	samplesCopy := make(map[string]Sample, len(c.samples))
	for id, sample := range c.samples {
		samplesCopy[id] = sample
	}
	return samplesCopy
}