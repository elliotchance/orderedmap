package orderedmap

import "sync"

type SyncOrderedMap[K comparable, V any] struct {
	OrderedMap[K, V]
	sync.RWMutex
}

func NewSyncOrderedMap[K comparable, V any]() *SyncOrderedMap[K, V] {
	return &SyncOrderedMap[K, V]{*NewOrderedMap[K, V](), sync.RWMutex{}}
}

func NewSyncOrderedMapWithCapacity[K comparable, V any](capacity int) *SyncOrderedMap[K, V] {
	return &SyncOrderedMap[K, V]{*NewOrderedMapWithCapacity[K, V](capacity), sync.RWMutex{}}
}

func (m *SyncOrderedMap[K, V]) Get(key K) (value V, ok bool) {
	m.RLock()
	defer m.RUnlock()

	return m.OrderedMap.Get(key)
}

func (m *SyncOrderedMap[K, V]) Set(key K, value V) bool {
	m.Lock()
	defer m.Unlock()

	return m.OrderedMap.Set(key, value)
}

func (m *SyncOrderedMap[K, V]) ReplaceKey(originalKey, newKey K) bool {
	m.Lock()
	defer m.Unlock()

	return m.OrderedMap.ReplaceKey(originalKey, newKey)
}

func (m *SyncOrderedMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	m.RLock()
	defer m.RUnlock()

	return m.OrderedMap.GetOrDefault(key, defaultValue)
}

func (m *SyncOrderedMap[K, V]) Len() int {
	m.RLock()
	defer m.RUnlock()

	return m.OrderedMap.Len()
}

func (m *SyncOrderedMap[K, V]) Delete(key K) (didDelete bool) {
	m.Lock()
	defer m.Unlock()

	return m.OrderedMap.Delete(key)
}

func (m *SyncOrderedMap[K, V]) Copy() *OrderedMap[K, V] {
	m.RLock()
	defer m.RUnlock()

	return m.OrderedMap.Copy()
}

func (m *SyncOrderedMap[K, V]) Has(key K) bool {
	m.RLock()
	defer m.RUnlock()

	return m.OrderedMap.Has(key)
}
