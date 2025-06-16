package mem

import (
	"container/list"
	"iter"
	"sync"
)

type SortedCache[K comparable, V any] struct {
	mu       sync.RWMutex
	mapCache map[K]*list.Element
	list     *list.List
}

type slot[K comparable, V any] struct {
	key K
	val V
}

func NewSortedCache[K comparable, V any]() *SortedCache[K, V] {
	return &SortedCache[K, V]{
		mapCache: make(map[K]*list.Element),
		list:     list.New(),
	}
}

func (cache *SortedCache[K, V]) Get(key K) (v V, ok bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	if elem, ok := cache.mapCache[key]; ok {
		return elem.Value.(slot[K, V]).val, true
	}
	return v, ok
}

func (cache *SortedCache[K, V]) Put(key K, value V) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if elem, ok := cache.mapCache[key]; ok {
		elem.Value = slot[K, V]{key, value}
		cache.list.MoveToBack(elem)
		return
	}
	cache.mapCache[key] = cache.list.PushBack(slot[K, V]{key, value})
}

func (cache *SortedCache[K, V]) Delete(key K) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if elem, ok := cache.mapCache[key]; ok {
		cache.list.Remove(elem)
		delete(cache.mapCache, elem.Value.(slot[K, V]).key)
	}
}

func (cache *SortedCache[K, V]) All() iter.Seq[V] {
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

func (cache *SortedCache[K, V]) Scan() iter.Seq2[K, V] {
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

func (cache *SortedCache[K, V]) Len() int {
	return len(cache.mapCache)
}
