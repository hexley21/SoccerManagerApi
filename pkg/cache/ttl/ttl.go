package ttl

import (
	"iter"
	"time"

	"github.com/hexley21/soccer-manager/pkg/cache"
)

type ExpirableItem[V any] struct {
	Value      V
	expiration time.Time
}

func NewItem[V any](value V, ttl time.Duration) ExpirableItem[V] {
	return ExpirableItem[V]{Value: value, expiration: time.Now().Add(ttl)}
}

type ttlCache[K comparable, V any] struct {
	inner cache.Cache[K, ExpirableItem[V]]
}

func New[K comparable, V any](
	inner cache.Cache[K, ExpirableItem[V]],
) cache.Cache[K, ExpirableItem[V]] {
	return &ttlCache[K, V]{inner: inner}
}

func (c *ttlCache[K, V]) Put(key K, item ExpirableItem[V]) {
	c.inner.Put(key, item)
}

func (c *ttlCache[K, V]) Get(key K) (ExpirableItem[V], bool) {
	item, ok := c.inner.Get(key)
	if !ok || time.Now().After(item.expiration) {
		c.inner.Delete(key)
		var zero ExpirableItem[V]
		return zero, false
	}
	return item, true
}

func (c *ttlCache[K, V]) Delete(key K) {
	c.inner.Delete(key)
}

func (c *ttlCache[K, V]) All() iter.Seq[ExpirableItem[V]] {
	return func(yield func(ExpirableItem[V]) bool) {
		now := time.Now()
		for item := range c.inner.All() {
			if now.After(item.expiration) {
				continue
			}
			if !yield(item) {
				return
			}
		}
	}
}

func (c *ttlCache[K, V]) Scan() iter.Seq2[K, ExpirableItem[V]] {
	return func(yield func(K, ExpirableItem[V]) bool) {
		now := time.Now()
		for k, item := range c.inner.Scan() {
			if now.After(item.expiration) {
				c.inner.Delete(k)
				continue
			}
			if !yield(k, item) {
				return
			}
		}
	}
}

func (c *ttlCache[K, V]) Len() int {
	return c.inner.Len()
}
