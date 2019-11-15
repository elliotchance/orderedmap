package orderedmap_test

import (
	"github.com/elliotchance/orderedmap"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElement_Key(t *testing.T) {
	t.Run("Front", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set(1, "foo")
		m.Set(2, "bar")
		assert.Equal(t, 1, m.Front().Key)
	})

	t.Run("Back", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set(1, "foo")
		m.Set(2, "bar")
		assert.Equal(t, 2, m.Back().Key)
	})
}

func TestElement_Value(t *testing.T) {
	t.Run("Front", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set(1, "foo")
		m.Set(2, "bar")
		assert.Equal(t, "foo", m.Front().Value)
	})

	t.Run("Back", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set(1, "foo")
		m.Set(2, "bar")
		assert.Equal(t, "bar", m.Back().Value)
	})
}

func TestElement_Next(t *testing.T) {
	m := orderedmap.NewOrderedMap()
	m.Set(1, "foo")
	m.Set(2, "bar")
	m.Set(3, "baz")

	var results []interface{}
	for el := m.Front(); el != nil; el = el.Next() {
		results = append(results, el.Key, el.Value)
	}

	assert.Equal(t, []interface{}{1, "foo", 2, "bar", 3, "baz"}, results)
}

func TestElement_Prev(t *testing.T) {
	m := orderedmap.NewOrderedMap()
	m.Set(1, "foo")
	m.Set(2, "bar")
	m.Set(3, "baz")

	var results []interface{}
	for el := m.Back(); el != nil; el = el.Prev() {
		results = append(results, el.Key, el.Value)
	}

	assert.Equal(t, []interface{}{3, "baz", 2, "bar", 1, "foo"}, results)
}
