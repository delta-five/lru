package lru

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	lock  sync.RWMutex
	pool  sync.Pool
	items map[K]*cacheItem[V]
	first *cacheItem[V]
	last  *cacheItem[V]
	len   int
	cap   int
}

type cacheItem[V any] struct {
	prev *cacheItem[V]
	next *cacheItem[V]
	mx   sync.Mutex
	val  V
}

func NewRLU[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		pool: sync.Pool{
			New: func() any { return new(cacheItem[V]) },
		},
		cap: capacity,
	}
}

func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.lock.RLock()
	item, ok := c.items[key]
	c.lock.RUnlock()
	if !ok {
		return
	}
	item.mx.Lock()
	defer item.mx.Unlock()

	return item.val, true
}
