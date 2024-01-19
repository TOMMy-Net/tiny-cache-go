package cache

import (
	"sync"
	"time"
)

type Cache struct {
    sync.RWMutex
    defaultExpiration time.Duration
    cleanupInterval   time.Duration
    items             map[string]Item
}

type Item struct {
    Value      interface{}
    Created    time.Time
    Expiration int64
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

    items := make(map[string]Item)

    cache := Cache{
        items:             items,
        defaultExpiration: defaultExpiration,
        cleanupInterval:   cleanupInterval,
    }

    // Если интервал очистки больше 0, запускаем GC (удаление устаревших элементов)
    if cleanupInterval > 0 {
       cache.StartGC() // данный метод рассматривается ниже
    }

    return &cache
}

func (c *Cache) StartGC()  {
	
}