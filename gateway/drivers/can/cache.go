package can

import "sync"

// Cache defines the interface for a local signal value cache.
type Cache interface {
	Set(value SignalValue)
	Get(id string) (SignalValue, bool)
	GetAll() []SignalValue
}

// MemoryCache is a thread-safe in-memory implementation of the Cache interface.
type MemoryCache struct {
	data map[string]SignalValue
	mu   sync.RWMutex
}

// NewMemoryCache creates a new in-memory cache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]SignalValue),
	}
}

// Set adds or updates a signal value in the cache.
func (c *MemoryCache) Set(value SignalValue) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[value.ID] = value
}

// Get retrieves a signal value from the cache by its ID.
func (c *MemoryCache) Get(id string) (SignalValue, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[id]
	return value, ok
}

// GetAll retrieves all signal values from the cache.
func (c *MemoryCache) GetAll() []SignalValue {
	c.mu.RLock()
	defer c.mu.RUnlock()
	values := make([]SignalValue, 0, len(c.data))
	for _, value := range c.data {
		values = append(values, value)
	}
	return values
}