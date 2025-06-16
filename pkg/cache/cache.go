package cache

import "iter"

type Cache[K comparable, V any] interface {
	Get(key K) (v V, ok bool)
	Put(key K, v V)
	Delete(key K)
	All() iter.Seq[V]
	Scan() iter.Seq2[K, V]
	Len() int
}
