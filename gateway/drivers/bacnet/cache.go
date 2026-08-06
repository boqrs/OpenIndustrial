package bacnet

import (
	"sync"
)

type Cache struct {
	mu sync.RWMutex

	samples map[string]Sample
}

func NewCache() *Cache {

	return &Cache{
		samples: make(map[string]Sample),
	}
}

// Update 更新一个 Sample
func (c *Cache) Update(sample Sample) {

	c.mu.Lock()
	defer c.mu.Unlock()

	c.samples[sample.ID] = sample
}

// BatchUpdate 批量更新
func (c *Cache) BatchUpdate(samples []Sample) {

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, sample := range samples {
		c.samples[sample.ID] = sample
	}
}

// Get 获取一个点位
func (c *Cache) Get(id string) (Sample, bool) {

	c.mu.RLock()
	defer c.mu.RUnlock()

	s, ok := c.samples[id]

	return s, ok
}

// Delete 删除
func (c *Cache) Delete(id string) {

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.samples, id)
}

// Clear 清空缓存
func (c *Cache) Clear() {

	c.mu.Lock()
	defer c.mu.Unlock()

	c.samples = make(map[string]Sample)
}

// Size 返回缓存数量
func (c *Cache) Size() int {

	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.samples)
}

// Snapshot 返回所有 Sample
func (c *Cache) Snapshot() []Sample {

	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]Sample, 0, len(c.samples))

	for _, s := range c.samples {
		result = append(result, s)
	}

	return result
}