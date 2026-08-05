package opcua

import "sync"

// Cache defines the interface for a local sample cache.
type Cache interface {
	Set(sample Sample)
	Get(id string) (Sample, bool)
	GetAll() []Sample
}

// MemoryCache is a thread-safe in-memory implementation of the Cache interface.
type MemoryCache struct {
	data map[string]Sample
	mu   sync.RWMutex
}

// NewMemoryCache creates a new in-memory cache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]Sample),
	}
}

// Set adds or updates a sample in the cache.
func (c *MemoryCache) Set(sample Sample) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[sample.ID] = sample
}

// Get retrieves a sample from the cache by its ID.
func (c *MemoryCache) Get(id string) (Sample, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	sample, ok := c.data[id]
	return sample, ok
}

// GetAll retrieves all samples from the cache.
func (c *MemoryCache) GetAll() []Sample {
	c.mu.RLock()
	defer c.mu.RUnlock()
	samples := make([]Sample, 0, len(c.data))
	for _, sample := range c.data {
		samples = append(samples, sample)
	}
	return samples
}