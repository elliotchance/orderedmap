//go:build go1.23
// +build go1.23

package orderedmap

import "iter"

func (m *OrderedMap[K, V]) Iterator() iter.Seq2[K, V] {
	return func(yield func(key K, value V) bool) {
		for el := m.Front(); el != nil; el = el.Next() {
			if !yield(el.Key, el.Value) {
				return
			}
		}
	}
}

func (m *OrderedMap[K, V]) ReverseIterator() iter.Seq2[K, V] {
	return func(yield func(key K, value V) bool) {
		for el := m.Back(); el != nil; el = el.Prev() {
			if !yield(el.Key, el.Value) {
				return
			}
		}
	}
}
