//go:build go1.23
// +build go1.23

package orderedmap_test

import (
	"testing"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/stretchr/testify/assert"
)

func TestIterators(t *testing.T) {
	type Element struct {
		Key   int
		Value bool
	}
	m := orderedmap.NewOrderedMap[int, bool]()
	expected := []Element{{5, true}, {3, false}, {1, false}, {4, true}}
	for _, v := range expected {
		m.Set(v.Key, v.Value)
	}

	t.Run("Iterator", func(t *testing.T) {
		i := 0
		for key, value := range m.Iterator() {
			assert.Equal(t, expected[i].Key, key)
			assert.Equal(t, expected[i].Value, value)
			i++
		}
	})

	t.Run("ReverseIterator", func(t *testing.T) {
		i := len(expected) - 1
		for key, value := range m.ReverseIterator() {
			assert.Equal(t, expected[i].Key, key)
			assert.Equal(t, expected[i].Value, value)
			i--
		}
	})
}
