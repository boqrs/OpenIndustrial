package ethercat

import "sync"

// Cache provides a thread-safe, in-memory storage for the latest samples.
// It's used to store the current state of all data points.
type Cache struct {
	mu     sync.RWMutex
	values map[string]Sample
}

// NewCache creates a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		values: make(map[string]Sample),
	}
}

// Set updates or adds a sample in the cache.
func (c *Cache) Set(sample Sample) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[sample.PointID] = sample
}

// Get retrieves a sample from the cache by its point ID.
func (c *Cache) Get(pointID string) (Sample, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	sample, found := c.values[pointID]
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