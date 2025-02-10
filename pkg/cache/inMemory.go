package cache

import (
	"sync"
	"time"
)

type CacheInMemory struct {
	storage    map[string]float64
	Ttl        int64
	LastUpdate int64
	sync.RWMutex
}

func NewCacheInMemory(ttl int64) *CacheInMemory {
	return &CacheInMemory{
		storage: make(map[string]float64),
		Ttl:     ttl,
	}
}

func (c *CacheInMemory) Set(new map[string]float64) {
	c.Lock()
	defer c.Unlock()

	c.storage = new
	c.LastUpdate = time.Now().Unix() + c.Ttl
}

func (c *CacheInMemory) Get() map[string]float64 {
	c.RLock()
	defer c.RUnlock()
	result := c.storage

	return result
}
