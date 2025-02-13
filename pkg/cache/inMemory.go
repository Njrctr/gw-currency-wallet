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

func (c *CacheInMemory) backgroundTask() {
	for {
		c.Lock()
		for key, el := range c.storage {
			if el.Ttl < time.Now().Unix() {
				delete(c.storage, key)
				fmt.Println("Удалён элемент ", el)
			}
		}
		c.Unlock()
		<-time.After(10 * time.Second)
	}
}

func NewCacheInMemory(ttl int64) *CacheInMemory {
	c := &CacheInMemory{
		storage: make(map[string]*InMemoryItem),
		ttl:     ttl,
	}

	go c.backgroundTask()

	return c
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
}

func (c *CacheInMemory) Get(key string) (*InMemoryItem, bool) {
	c.RLock()
	defer c.RUnlock()
	result, ex := c.storage[key]

	return result, ex
}
