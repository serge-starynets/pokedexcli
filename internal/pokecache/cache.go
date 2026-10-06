package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cacheMap map[string]cacheEntry
	mtx      sync.Mutex
	interval time.Duration
}

func (c *Cache) reap(now time.Time) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	for key, entry := range c.cacheMap {
		if now.Sub(entry.createdAt) > c.interval {
			delete(c.cacheMap, key)
		}
	}
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for now := range ticker.C {
		c.reap(now)
	}
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cacheMap: make(map[string]cacheEntry),
		interval: interval,
	}

	go c.reapLoop()

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.cacheMap[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	entry, ok := c.cacheMap[key]
	return entry.val, ok
}
