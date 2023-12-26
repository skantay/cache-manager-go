// Cache manager
// caman
// ca - cache
// man - manager
package caman

import (
	"encoding/json"
	"errors"
	"os"
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

func (c *Cache) Count() int {
	c.RLock()
	defer c.RUnlock()

	return len(c.items)
}

func (c *Cache) Rename(oldKey, newKey string) error {
	c.Lock()
	defer c.Unlock()

	if _, found := c.items[oldKey]; !found {
		return errors.New("Key not found")
	}

	item := c.items[oldKey]

	delete(c.items, oldKey)

	c.items[newKey] = item

	return nil
}

func (c *Cache) FlushAll() {
	c.Lock()
	defer c.Unlock()

	for key := range c.items {
		delete(c.items, key)
	}
}

func (c *Cache) Exist(value interface{}) bool {
	c.RLock()
	defer c.RUnlock()

	for _, v := range c.items {
		if v.Value == value {
			return true
		}
	}

	return false
}

func (c *Cache) Copy(key string) error {
	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]
	if !found {
		return errors.New("Key not found")
	}

	newKey := key + "_copy"
	c.items[newKey] = Item{
		Value:      item.Value,
		Expiration: item.Expiration,
		Created:    time.Now(),
	}

	return nil
}

func (c *Cache) Increment(key string, delta int) (int, error) {
	c.Lock()
	defer c.Unlock()

	item, found := c.items[key]
	if !found {
		return 0, errors.New("Key not found")
	}

	if intValue, ok := item.Value.(int); ok {
		newValue := intValue + delta
		c.items[key] = Item{
			Value:      newValue,
			Expiration: item.Expiration,
			Created:    time.Now(),
		}
		return newValue, nil
	}

	return 0, errors.New("Value is not an integer")
}

func (c *Cache) Decrement(key string, delta int) (int, error) {
	return c.Increment(key, -delta)
}

func (c *Cache) Expire(key string) bool {
	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	if !found {
		return true
	}

	if item.Expiration > 0 {

		if time.Now().UnixNano() > item.Expiration {
			return true
		}
	}

	return false
}

func (c *Cache) SaveFileJSON(filename string) error {
	c.RLock()
	defer c.RUnlock()

	// Serialize cache to JSON
	serialized, err := json.Marshal(c.items)
	if err != nil {
		return err
	}

	// Write JSON data to the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(serialized)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) LoadFileJSON(filename string) error {
    // Read JSON data from the file
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    // Decode JSON data into map[string]Item
    var items map[string]Item
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&items); err != nil {
        return err
    }

    c.Lock()
    defer c.Unlock()

    // Replace the existing cache items with the loaded items
    c.items = items

    return nil
}