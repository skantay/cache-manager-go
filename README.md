Certainly! Below is a small README template for your Go caching project:

---

# Caman: Simple Go Cache Manager

Caman is a lightweight cache manager for Go that provides easy-to-use caching functionalities. It allows you to store and retrieve data in-memory, with options for automatic expiration, key renaming, and more.

## Features

- **Set:** Store a key-value pair with an optional expiration duration.
- **Get:** Retrieve a value by its key, with automatic expiration check.
- **Delete:** Remove a key from the cache.
- **Count:** Get the current count of items in the cache.
- **Rename:** Change the name of a key in the cache.
- **FlushAll:** Clear all items from the cache.
- **Exist:** Check if a value exists in the cache.
- **Copy** Creates a duplicate of an item
- **Increment** Increments an int type value
- **Decrement** Decrements an int type value
- **Expire** Check if an item is expired
- **SaveFileJSON** Save all items into a json file
- **LoadFileJSON** Get all items from a json file

## Usage

```go
import "github.com/skantay/caman"

// Example usage
cache := caman.New(defaultExpiration, cleanupInterval)

// Set a value with expiration
cache.Set("key", "value", time.Minute)

// Get a value
value, found := cache.Get("key")

// Delete a key
cache.Delete("key")

// Rename a key
cache.Rename("oldKey", "newKey")

// Check if a value exists
exists := cache.Exist("value")
```

## Serialization

Caman supports serialization to and from files using JSON.

### JSON Serialization

```go
// Save cache to file
cache.SaveFileJSON("cache.json")

// Load cache from file
cache.LoadFileJSON("cache.json")
```