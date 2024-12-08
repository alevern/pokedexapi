package cache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type cacheMap struct {
	cache    map[string]cacheEntry
	interval time.Duration
	mu       sync.RWMutex
}

type Cache interface {
	Add(key string, val []byte)
	Get(key string) ([]byte, bool)
	deleteOldValues()
	reapLoop()
}

func NewCache(interval time.Duration) Cache {
	cache := &cacheMap{
		cache:    make(map[string]cacheEntry),
		interval: interval,
	}
	cache.reapLoop()
	return cache
}

func (c *cacheMap) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// if _, prs := c.cache[key]; prs {
	// 	return fmt.Errorf("[ADD] key %s already exists", key)
	// }
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	// return nil
}

func (c *cacheMap) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cacheEntry, prs := c.cache[key]
	if prs == false {
		return nil, false
	}
	return cacheEntry.val, true
}

func (c *cacheMap) reapLoop() {
	ticker := time.NewTicker(c.interval)
	go func() {
		for range ticker.C {
			c.deleteOldValues()
		}
	}()
}

func (c *cacheMap) deleteOldValues() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, val := range c.cache {
		if time.Now().Sub(val.createdAt) > c.interval {
			delete(c.cache, key)
		}
	}
}
