// Cache manager
// caman
// ca - cache
// man - manager
package caman

import (
	"errors"
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

	// Initializing map
	items := make(map[string]Item)

	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	// If cleanUp interval is greater than 0, then start GC(process of deleting old elements)
	if cleanupInterval > 0 {
		cache.startGC()
	}

	return &cache
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	if !found {
		return nil, false
	}

	if item.Expiration > 0 {

		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}

	return item.Value, true
}

func (c *Cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("Key not found")
	}

	delete(c.items, key)

	return nil
}

func (c *Cache) startGC() {
	go c.GC()
}

func (c *Cache) GC() {
	for {
		<- time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}


func (c *Cache) expiration() (keys []string) {
	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().Unix() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return 
}

func (c *Cache) clearItems(keys []string) {
	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}

// TO-DO

/*
Count — получение кол-ва элементов в кеше
GetItem — получение элемента кеша
Rename — переименования ключа
Copy — копирование элемента
Increment — инкремент
Decrement — декремент
Exist — проверка элемента на существование
Expire — проверка кеша на истечение срока жизни
FlushAll — очистка всех данных
SaveFile — сохранение данных в файл
LoadFile — загрузка данных из файла


*/