package mem

import (
	"iter"
	"sync"
)

type InMemoryCache[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func NewInMemoryCache[K comparable, V any]() *InMemoryCache[K, V] {
	return &InMemoryCache[K, V]{
		data: make(map[K]V),
	}
}

func (c *InMemoryCache[K, V]) Get(k K) (v V, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok = c.data[k]
	return
}

func (c *InMemoryCache[K, V]) Put(k K, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[k] = v
}

func (c *InMemoryCache[K, V]) Delete(k K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, k)
}

func (c *InMemoryCache[K, V]) All() iter.Seq[V] {
	return func(yield func(V) bool) {
		c.mu.RLock()
		defer c.mu.RUnlock()

		for _, v := range c.data {
			if !yield(v) {
				return
			}
		}
	}
}

func (c *InMemoryCache[K, V]) Scan() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		c.mu.RLock()
		defer c.mu.RUnlock()

		for k, v := range c.data {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (c *InMemoryCache[K, V]) Len() int {
	return len(c.data)
}
