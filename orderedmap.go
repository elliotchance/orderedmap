package orderedmap

import "sync"

type OrderedMap struct {
	sync.RWMutex
	kv map[interface{}]*Element
	ll list
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		kv: make(map[interface{}]*Element),
	}
}

// NewOrderedMapWithCapacity creates a map with enough pre-allocated space to
// hold the specified number of elements.
func NewOrderedMapWithCapacity(capacity int) *OrderedMap {
	return &OrderedMap{
		kv: make(map[interface{}]*Element, capacity),
	}
}

// Get returns the value for a key. If the key does not exist, the second return
// parameter will be false and the value will be nil.
func (m *OrderedMap) Get(key interface{}) (interface{}, bool) {
	m.RLock()
	defer m.RUnlock()

	element, ok := m.kv[key]
	if ok {
		return element.Value, true
	}

	return nil, false
}

// Set will set (or replace) a value for a key. If the key was new, then true
// will be returned. The returned value will be false if the value was replaced
// (even if the value was the same).
func (m *OrderedMap) Set(key, value interface{}) bool {
	m.Lock()
	defer m.Unlock()

	_, alreadyExist := m.kv[key]
	if alreadyExist {
		m.kv[key].Value = value
		return false
	}

	element := m.ll.PushBack(key, value)
	m.kv[key] = element
	return true
}

// GetOrDefault returns the value for a key. If the key does not exist, returns
// the default value instead.
func (m *OrderedMap) GetOrDefault(key, defaultValue interface{}) interface{} {
	m.RLock()
	defer m.RUnlock()

	if element, ok := m.kv[key]; ok {
		return element.Value
	}

	return defaultValue
}

// GetElement returns the element for a key. If the key does not exist, the
// pointer will be nil.
func (m *OrderedMap) GetElement(key interface{}) *Element {
	m.RLock()
	defer m.RUnlock()

	element, ok := m.kv[key]
	if ok {
		return element
	}

	return nil
}

// Len returns the number of elements in the map.
func (m *OrderedMap) Len() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.kv)
}

// Keys returns all of the keys in the order they were inserted. If a key was
// replaced it will retain the same position. To ensure most recently set keys
// are always at the end you must always Delete before Set.
func (m *OrderedMap) Keys() (keys []interface{}) {
	m.RLock()
	defer m.RUnlock()

	keys = make([]interface{}, 0, m.Len())
	for el := m.Front(); el != nil; el = el.Next() {
		keys = append(keys, el.Key)
	}
	return keys
}

// Delete will remove a key from the map. It will return true if the key was
// removed (the key did exist).
func (m *OrderedMap) Delete(key interface{}) (didDelete bool) {
	m.Lock()
	defer m.Unlock()

	element, ok := m.kv[key]
	if ok {
		m.ll.Remove(element)
		delete(m.kv, key)
	}

	return ok
}

// Front will return the element that is the first (oldest Set element). If
// there are no elements this will return nil.
func (m *OrderedMap) Front() *Element {
	m.RLock()
	defer m.RUnlock()

	return m.ll.Front()
}

// Back will return the element that is the last (most recent Set element). If
// there are no elements this will return nil.
func (m *OrderedMap) Back() *Element {
	m.RLock()
	defer m.RUnlock()

	return m.ll.Back()
}

// Copy returns a new OrderedMap with the same elements.
// Using Copy while there are concurrent writes may mangle the result.
func (m *OrderedMap) Copy() *OrderedMap {
	m.RLock()
	defer m.RUnlock()

	m2 := NewOrderedMapWithCapacity(m.Len())

	for el := m.Front(); el != nil; el = el.Next() {
		m2.Set(el.Key, el.Value)
	}

	return m2
}

// Has checks if a key exists in the map.
func (m *OrderedMap) Has(key interface{}) bool {
	m.RLock()
	defer m.RUnlock()

	_, exists := m.kv[key]
	return exists
}
