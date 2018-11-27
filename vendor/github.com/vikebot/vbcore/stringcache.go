package vbcore

import (
	"sync"
	"time"
)

// StringCacheRefreshFunc defines a func that can used to update the value of an Item
type StringCacheRefreshFunc func() (value string)

// StringCacheItem is one single object of the Cache
type StringCacheItem struct {
	sync.RWMutex
	last     time.Time
	value    string
	refresh  StringCacheRefreshFunc
	interval float64
}

// StringCache defines a simple to use cache for string values
type StringCache struct {
	storage map[string]*StringCacheItem
}

// NewStringCache creates the internally used storages
func NewStringCache() *StringCache {
	return &StringCache{
		storage: make(map[string]*StringCacheItem),
	}
}

// Add registers a new cache item that can be received with Get(). All cache
// items will automatically be updated in the specified intervals. Interval
// values should be in second unity
func (c *StringCache) Add(key string, intervalSeconds float64, refresh StringCacheRefreshFunc) {
	c.storage[key] = &StringCacheItem{
		interval: intervalSeconds,
		refresh:  refresh,
	}
	c.Get(key)
}

// Get returns the last fetched value of the cache item. If in this Get cycle
// the cache detects that an refresh is needed you will get the old value and a
// refresh will be executed in another go rountine.
func (c *StringCache) Get(key string) (value string) {
	cs := c.storage[key]
	cs.RLock()
	// Check if the value is up-to-date
	if cs.last.Year() == 1 || time.Now().Sub(cs.last).Seconds() >= cs.interval {
		// Run update in new go rountine so the current user doesn't suffer
		// from performance
		go func() {
			updatedVal := cs.refresh()
			cs.Lock()
			cs.value = updatedVal
			cs.last = time.Now()
			cs.Unlock()
		}()
	}
	value = cs.value
	cs.RUnlock()
	return
}
