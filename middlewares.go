package tinydb

import (
	"encoding/json"
)

// MiddlewareCaching is a middleware for caching data.
type MiddlewareCaching struct {
	storage Storage
	cache   []byte
	count   int
	size    int
}

// CachingMiddleware Create a new CachingMiddleware instance.
func CachingMiddleware(storage Storage, size int) *MiddlewareCaching {
	return &MiddlewareCaching{storage: storage, size: size}
}

// Read data from MiddlewareCaching cache.
func (mv *MiddlewareCaching) Read(data any) error {
	if mv.cache == nil {
		return mv.storage.Read(data)
	}
	return json.Unmarshal(mv.cache, &data)
}

// Write data to MiddlewareCaching cache.
func (mv *MiddlewareCaching) Write(data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	mv.cache = b
	mv.count++
	if mv.count >= mv.size {
		err := mv.Read(data)
		if err != nil {
			return err
		}
		return mv.storage.Write(data)
	}
	return nil
}

// Close the MemoryStorage instance.
func (mv *MiddlewareCaching) Close() error {
	var data = make(map[string]interface{})
	err := mv.Read(data)
	if err != nil {
		return err
	}
	err = mv.storage.Write(data)
	if err != nil {
		return err
	}
	return mv.storage.Close()
}
