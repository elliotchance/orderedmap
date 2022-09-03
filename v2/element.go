package orderedmap

import (
	"container/list"
)

type Element[K comparable, V any] struct {
	Key   K
	Value V

	element *list.Element
}

func newElement[K comparable, V any](e *list.Element) *Element[K, V] {
	if e == nil {
		return nil
	}

	element := e.Value.(*orderedMapElement[K, V])

	return &Element[K, V]{
		element: e,
		Key:     element.key,
		Value:   element.value,
	}
}

// Next returns the next element, or nil if it finished.
func (e *Element[K, V]) Next() *Element[K, V] {
	return newElement[K, V](e.element.Next())
}

// Prev returns the previous element, or nil if it finished.
func (e *Element[K, V]) Prev() *Element[K, V] {
	return newElement[K, V](e.element.Prev())
}
