package utils

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type InMemoryCache struct {
	data map[string]CacheItem
	mu   sync.RWMutex
	ttl time.Duration
}

func NewInMemoryCache(ttl time.Duration) *InMemoryCache{
	return  &InMemoryCache{
		data: make(map[string]CacheItem),
		ttl: ttl,
	}
}

func (c *InMemoryCache) Set(key string, value interface{}){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] =CacheItem{
		Value: value,
		Expiration: time.Now().Add(c.ttl).Unix(),
	}
}

func (c *InMemoryCache) Get(key string) (interface{}, bool){
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.data[key]
	if !found || time.Now().Unix()> item.Expiration{
		return nil, false
	}
	return item.Value, true
}

func (c *InMemoryCache) Delete(key string){
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}