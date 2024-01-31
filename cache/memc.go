package cache

import (
	"errors"
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
type Cleaner interface {
	DealeteEx()
}

type Result interface {
	String() (string)
	Byte() ([]byte, error)
	Int() (int, error)
	Float64() (float64, error)
}



const (
	NotByte    = "This type not []byte"
	NotInt     = "This type not int"
	NotFloat64 = "This type not float64"
	NotFound = "This key not found in cache"
)

const (
	defaultExpirationConst      = 720 * time.Hour  // 1 month
	defaultcleanupIntervalConst = 10 * time.Second // 10 sec
	defaultBuf                  = 500
)

// Creating a new cache space
func New() *Cache {

	items := make(map[string]Item, defaultBuf)

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
func (c *Cache) Get(key string) Result {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.Items[key]
	if !found {
		return Item{}
	} else {
		if item.Expiration > time.Now().UnixMilli() {
			return item

		}
	}
	return Item{}

}

// Return key value and delete key
func (c *Cache) GetD(key string) Result {
	rst := c.Get(key)
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Items, key)
	return rst
}

// Always return string
func (i Item) String() (string) {
	if i.Value != nil {
		return fmt.Sprint(i.Value)
	}
	return ""
}

// Return []byte
func (i Item) Byte() ([]byte, error) {
	if i.Value != nil {
		if v, ok := i.Value.([]byte); ok {
			return v, nil
		} else {
			return nil, errors.New(NotByte)
		}
	}
	return nil, nil
}

// Return int
func (i Item) Int() (int, error) {
	if i.Value != nil {
		switch i.Value.(type) {
		case float32:
			if v, ok := i.Value.(float32); ok {

				return int(v), nil
			}
		case float64:
			if v, ok := i.Value.(float64); ok {

				return int(v), nil
			}
		case int:
			if v, ok := i.Value.(int); ok {

				return v, nil
			}
		case int8:
			if v, ok := i.Value.(int8); ok {

				return int(v), nil
			}
		case int16:
			if v, ok := i.Value.(int16); ok {

				return int(v), nil
			}
		case int32:
			if v, ok := i.Value.(int32); ok {

				return int(v), nil
			}
		case int64:
			if v, ok := i.Value.(int64); ok {

				return int(v), nil
			}
		default:
			return 0, errors.New(NotInt)
		}
	}
	return 0, errors.New(NotFound)
}

func (i Item) Float64() (float64, error) {
	if i.Value != nil {
		switch i.Value.(type) {
		case float32:
			if v, ok := i.Value.(float32); ok {

				return float64(v), nil
			}
		case float64:
			if v, ok := i.Value.(float64); ok {
				return v, nil
			}
		case int:
			if v, ok := i.Value.(int); ok {

				return float64(v), nil
			}
		case int8:
			if v, ok := i.Value.(int8); ok {

				return float64(v), nil
			}
		case int16:
			if v, ok := i.Value.(int16); ok {

				return float64(v), nil
			}
		case int32:
			if v, ok := i.Value.(int32); ok {

				return float64(v), nil
			}
		case int64:
			if v, ok := i.Value.(int64); ok {

				return float64(v), nil
			}
		default:
			return 0, errors.New(NotFloat64)
		}

	}
	return 0, errors.New(NotFound)
}

// Return full cache
func (c *Cache) GetFullMap() map[string]Item {
	return c.Items
}

// Return len of data
func (c *Cache) Count() int {
	return len(c.Items)
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

	new := make(map[string]Item, defaultBuf)
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
