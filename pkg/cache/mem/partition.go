package mem

import (
	"fmt"
	"hash/fnv"
	"iter"
)

type PartitionCache[K comparable, V any] struct {
	partitions []InMemoryCache[K, V]
	size       int
}

func NewPartitionCache[K comparable, V any](n int) *PartitionCache[K, V] {
	partitions := make([]InMemoryCache[K, V], n)
	for i := range n {
		partitions[i] = InMemoryCache[K, V]{data: make(map[K]V)}
	}
	return &PartitionCache[K, V]{
		partitions: partitions,
		size:       n,
	}
}

func (c *PartitionCache[K, V]) hash(k K) int {
	h := fnv.New32a()
	fmt.Fprintf(h, "%v", k)
	return int(h.Sum32()) % c.size
}

func (c *PartitionCache[K, V]) Get(k K) (v V, ok bool) {
	i := c.hash(k)
	c.partitions[i].mu.RLock()
	defer c.partitions[i].mu.RUnlock()
	v, ok = c.partitions[i].data[k]
	return
}

func (c *PartitionCache[K, V]) Put(key K, v V) {
	i := c.hash(key)
	c.partitions[i].mu.Lock()
	defer c.partitions[i].mu.Unlock()
	c.partitions[i].data[key] = v
}

func (c *PartitionCache[K, V]) Delete(key K) {
	i := c.hash(key)
	c.partitions[i].mu.Lock()
	defer c.partitions[i].mu.Unlock()
	delete(c.partitions[i].data, key)
}

func (c *PartitionCache[K, V]) All() iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := range c.size {
			c.partitions[i].mu.RLock()
			defer c.partitions[i].mu.RUnlock()

			for _, v := range c.partitions[i].data {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func (c *PartitionCache[K, V]) Scan() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i := range c.size {
			c.partitions[i].mu.RLock()
			defer c.partitions[i].mu.RUnlock()

			for k, v := range c.partitions[i].data {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

func (c *PartitionCache[K, V]) Len() int {
	var l int
	for i := range c.size {
		l += len(c.partitions[i].data)
	}

	return l
}
