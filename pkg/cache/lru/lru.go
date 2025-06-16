package lru

import (
	"container/list"
	"iter"
	"sync"
)

type LRUCache[K comparable, V any] struct {
	mu       sync.RWMutex
	mapCache map[K]*list.Element
	list     *list.List
	capacity int
}

type slot[K comparable, V any] struct {
	key K
	val V
}

func New[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		mapCache: make(map[K]*list.Element, capacity),
		list:     list.New(),
		capacity: capacity,
	}
}

func (cache *LRUCache[K, V]) Get(key K) (v V, ok bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	if elem, ok := cache.mapCache[key]; ok {
		cache.list.MoveToFront(elem)
		return elem.Value.(slot[K, V]).val, true
	}
	return v, ok
}

func (cache *LRUCache[K, V]) Put(key K, value V) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if elem, ok := cache.mapCache[key]; ok {
		elem.Value = slot[K, V]{key, value}
		cache.list.MoveToFront(elem)
		return
	}
	if len(cache.mapCache) == cache.capacity {
		tail := cache.list.Back()
		delete(cache.mapCache, tail.Value.(slot[K, V]).key)
		cache.list.Remove(tail)
	}
	cache.mapCache[key] = cache.list.PushFront(slot[K, V]{key, value})
}

func (cache *LRUCache[K, V]) Delete(key K) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if elem, ok := cache.mapCache[key]; ok {
		cache.list.Remove(elem)
		delete(cache.mapCache, elem.Value.(slot[K, V]).key)
	}
}

func (cache *LRUCache[K, V]) All() iter.Seq[V] {
	return func(yield func(V) bool) {
		cache.mu.RLock()
		defer cache.mu.RUnlock()

		for e := cache.list.Front(); e != nil; e = e.Next() {
			item := e.Value.(slot[K, V])
			if !yield(item.val) {
				return
			}
		}
	}
}

func (cache *LRUCache[K, V]) Scan() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		cache.mu.RLock()
		defer cache.mu.RUnlock()

		for e := cache.list.Front(); e != nil; e = e.Next() {
			item := e.Value.(slot[K, V])
			if !yield(item.key, item.val) {
				return
			}
		}
	}
}

func (cache *LRUCache[K, V]) Len() int {
	return len(cache.mapCache)
}
