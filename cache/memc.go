package cache

import (
	"fmt"
	"sync"
	"time"
)

// Cache is a simple caching system that stores key-value pairs.
type Cache struct {
	mu                sync.RWMutex
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
	Items             map[string]Item
}

// Struct of the vallue
type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

const (
	defaultExpirationConst      = 5 * time.Hour    // 6 hour
	defaultcleanupIntervalConst = 10 * time.Second // 10 sec
	defaultBuf                  = 200
)

// Creating a new cache space
func New() *Cache {

	items := make(map[string]Item)

	cache := Cache{
		Items:             items,
		DefaultExpiration: defaultExpirationConst,
		CleanupInterval:   defaultcleanupIntervalConst,
	}

	go StartGC(&cache)

	return &cache
}

// Set a new default cache expiration time
func (c *Cache) SetDefaultExpiration(t time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.DefaultExpiration = t
}


// Set a new default cache CleanupInterval time
func (c *Cache) SetDefaultCleanupInterval(t time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CleanupInterval = t
}

// Add an item to the cache, replacing any existing item.
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	if duration == 0 || duration < 0 {
		duration = c.DefaultExpiration
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Items[key] = Item{
		Value:      value,
		Created:    time.Now(),
		Expiration: time.Now().Add(duration).UnixMilli(),
	}
}

// Retern key vallue
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.Items[key]
	if !found {
		return "", false
	} else {
		if item.Expiration > time.Now().UnixMilli() {
			return fmt.Sprint(item.Value), true
		}
	}
	return "", false

}

// Return full cache
func (c *Cache) GetFullMap() map[string]Item {
	return c.Items
}

// Returns the key expiration time in Unix
func (c *Cache) GetExUnix(key string) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if v, ok := c.Items[key]; ok {
		return v.Expiration
	}
	return 0
}

func StartGC(c *Cache) {
	for {
		<-time.After(c.CleanupInterval)
		if c.Items != nil {
			c.DealeteEx()
		}
	}
}

// Dealete all in cache
func (c *Cache) DealeteAllCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	new := make(map[string]Item)
	c.Items = new
}


func (c *Cache) DealeteEx() {
	c.mu.Lock()

	new := make(map[string]Item)
	for i, k := range c.Items {
		if time.Now().UnixMilli() < k.Expiration && k.Expiration > 0 {
			new[i] = k
		}
	}
	c.Items = new
	c.mu.Unlock()
}

func AddNewMap(c *Cache, n map[string]Item) {
	if n != nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.Items = n
	}
}
