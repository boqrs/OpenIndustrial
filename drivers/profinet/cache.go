package profinet

import "sync"

// Cache provides a thread-safe in-memory storage for PROFINET samples.
type Cache struct {
	mu     sync.RWMutex
	values map[string]Sample
}

// NewCache creates a new instance of a Cache.
func NewCache() *Cache {
	return &Cache{
		values: make(map[string]Sample),
	}
}

// Update adds or updates a sample in the cache.
func (c *Cache) Update(sample Sample) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[sample.PointID] = sample
}

// Get retrieves a sample from the cache by its PointID.
func (c *Cache) Get(id string) (Sample, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	sample, found := c.values[id]
	return sample, found
}

// GetAll returns a copy of all samples currently in the cache.
func (c *Cache) GetAll() []Sample {
	c.mu.RLock()
	defer c.mu.RUnlock()
	samples := make([]Sample, 0, len(c.values))
	for _, sample := range c.values {
		samples = append(samples, sample)
	}
	return samples
}