package orderedmap

import "sync"

type SyncOrderedMap[K comparable, V any] struct {
	OrderedMap[K, V]
	sync.RWMutex
}
