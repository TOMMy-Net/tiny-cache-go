package cache

import (
	"fmt"

	"sync"
	"time"
)

type Cache struct {
	mu                sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item
}

type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

const (
	defaultExpirationConst = 6 * time.Hour // 6 hour
	defaultcleanupIntervalConst   = 10 * time.Second // 10 sec
)

func New() *Cache {

	items := make(map[string]Item)

	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpirationConst,
		cleanupInterval:   defaultcleanupIntervalConst,
	}

	go StartGC(&cache)

	return &cache
}

func (c *Cache) SetDefaultExpiration(t time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.defaultExpiration = t
}

func (c *Cache) SetDefaultCleanupInterval(t time.Duration)  {
    c.mu.Lock()
	defer c.mu.Unlock()
	c.cleanupInterval = t
}

// Add an item to the cache, replacing any existing item.
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	if duration == 0 || duration < 0 {
		duration = c.defaultExpiration
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = Item{
		Value:      value,
		Created:    time.Now(),
		Expiration: time.Now().Add(duration).UnixMilli(),
	}
}
func (c *Cache) Get(k string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[k]
	if !found {
		return "", false
	} else {
		if item.Expiration > time.Now().UnixMilli() {
			return fmt.Sprint(item.Value), true
		}
	}
	return "", false

}
func (c *Cache) GetFullMap() map[string]Item {
	return c.items
}

func StartGC(c *Cache) {
	for {
		<-time.After(c.cleanupInterval)
		if c.items != nil {
			c.DealeteEx()
		}
	}
}

func (c *Cache) DealeteEx() {
	c.mu.Lock()
	
	new := make(map[string]Item)
	for i, k := range c.items {
		if time.Now().UnixMilli() < k.Expiration && k.Expiration > 0 {
			new[i] = k
		}
	}
	c.items = new
	c.mu.Unlock()
}
 
func AddNewMap(c *Cache, n map[string]Item)  {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = n
}