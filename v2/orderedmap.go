package orderedmap

import (
	"container/list"
	"golang.org/x/exp/constraints"
)

type orderedMapElement[K constraints.Ordered, V any] struct {
	key   K
	value V
}

type OrderedMap[K constraints.Ordered, V any] struct {
	kv map[K]*list.Element
	ll *list.List
}

func NewOrderedMap[K constraints.Ordered, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		kv: make(map[K]*list.Element),
		ll: list.New(),
	}
}

// Get returns the value for a key. If the key does not exist, the second return
// parameter will be false and the value will be nil.
func (m *OrderedMap[K, V]) Get(key K) (value V, ok bool) {
	v, ok := m.kv[key]
	if ok {
		value = v.Value.(*orderedMapElement[K, V]).value
	}

	return
}

// Set will set (or replace) a value for a key. If the key was new, then true
// will be returned. The returned value will be false if the value was replaced
// (even if the value was the same).
func (m *OrderedMap[K, V]) Set(key K, value V) bool {
	_, didExist := m.kv[key]

	if !didExist {
		element := m.ll.PushBack(&orderedMapElement[K, V]{key, value})
		m.kv[key] = element
	} else {
		m.kv[key].Value.(*orderedMapElement[K, V]).value = value
	}

	return !didExist
}

// GetOrDefault returns the value for a key. If the key does not exist, returns
// the default value instead.
func (m *OrderedMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	if value, ok := m.kv[key]; ok {
		return value.Value.(*orderedMapElement[K, V]).value
	}

	return defaultValue
}

// GetElement returns the element for a key. If the key does not exist, the
// pointer will be nil.
func (m *OrderedMap[K, V]) GetElement(key K) *Element[K, V] {
	value, ok := m.kv[key]
	if ok {
		element := value.Value.(*orderedMapElement[K, V])
		return &Element[K, V]{
			element: value,
			Key:     element.key,
			Value:   element.value,
		}
	}

	return nil
}

// Len returns the number of elements in the map.
func (m *OrderedMap[K, V]) Len() int {
	return len(m.kv)
}

// Keys returns all of the keys in the order they were inserted. If a key was
// replaced it will retain the same position. To ensure most recently set keys
// are always at the end you must always Delete before Set.
func (m *OrderedMap[K, V]) Keys() (keys []K) {
	keys = make([]K, m.Len())

	element := m.ll.Front()
	for i := 0; element != nil; i++ {
		keys[i] = element.Value.(*orderedMapElement[K, V]).key
		element = element.Next()
	}

	return keys
}

// Delete will remove a key from the map. It will return true if the key was
// removed (the key did exist).
func (m *OrderedMap[K, V]) Delete(key K) (didDelete bool) {
	element, ok := m.kv[key]
	if ok {
		m.ll.Remove(element)
		delete(m.kv, key)
	}

	return ok
}

// Front will return the element that is the first (oldest Set element). If
// there are no elements this will return nil.
func (m *OrderedMap[K, V]) Front() *Element[K, V] {
	front := m.ll.Front()
	if front == nil {
		return nil
	}

	element := front.Value.(*orderedMapElement[K, V])

	return &Element[K, V]{
		element: front,
		Key:     element.key,
		Value:   element.value,
	}
}

// Back will return the element that is the last (most recent Set element). If
// there are no elements this will return nil.
func (m *OrderedMap[K, V]) Back() *Element[K, V] {
	back := m.ll.Back()
	if back == nil {
		return nil
	}

	element := back.Value.(*orderedMapElement[K, V])

	return &Element[K, V]{
		element: back,
		Key:     element.key,
		Value:   element.value,
	}
}

// Copy returns a new OrderedMap with the same elements.
// Using Copy while there are concurrent writes may mangle the result.
func (m *OrderedMap[K, V]) Copy() *OrderedMap[K, V] {
	m2 := NewOrderedMap[K, V]()

	for el := m.Front(); el != nil; el = el.Next() {
		m2.Set(el.Key, el.Value)
	}

	return m2
}
