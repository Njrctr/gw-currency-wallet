package cache

import (
	"fmt"
	"sync"
	"time"
)

type CacheInMemory struct {
	storage map[string]*InMemoryItem
	ttl     int64
	sync.RWMutex
}

func NewCacheInMemory(ttl int64) *CacheInMemory {
	return &CacheInMemory{
		storage: make(map[string]*InMemoryItem),
	}
}

type InMemoryItem struct {
	Value float64
	Ttl   int64
}

func (c *CacheInMemory) Set(key string, val float64) {
	c.Lock()
	defer c.Unlock()

	c.storage[key] = &InMemoryItem{
		Value: val,
		Ttl:   time.Now().Unix() + c.ttl,
	}
	fmt.Println(c.storage)
}

func (c *CacheInMemory) Get(key string) *InMemoryItem {
	c.RLock()
	defer c.RUnlock()
	result := c.storage[key]

	return result
}
