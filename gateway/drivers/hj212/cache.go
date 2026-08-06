package hj212

import (
	"sync"
)

// Cache provides a thread-safe in-memory store for the latest sample data.
type Cache struct {
	mu    sync.RWMutex
	items map[string]Sample
}

// NewCache creates a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]Sample),
	}
}

// Set adds or updates a sample in the cache.
func (c *Cache) Set(pointID string, sample Sample) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[pointID] = sample
}

// Get retrieves a sample from the cache by its point ID.
func (c *Cache) Get(pointID string) (Sample, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	sample, found := c.items[pointID]
	return sample, found
}

// GetAll returns a copy of all items in the cache.
func (c *Cache) GetAll() map[string]Sample {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	// Create a copy to prevent race conditions on the returned map.
	allItems := make(map[string]Sample, len(c.items))
	for key, value := range c.items {
		allItems[key] = value
	}
	return allItems
}