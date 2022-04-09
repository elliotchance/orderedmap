package orderedmap_test

import (
	"testing"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderedMap(t *testing.T) {
	m := orderedmap.NewOrderedMap[int, string]()
	assert.IsType(t, &orderedmap.OrderedMap[int, string]{}, m)
}

func TestGet(t *testing.T) {
	t.Run("ReturnsNotOKIfStringKeyDoesntExist", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		_, ok := m.Get("foo")
		assert.False(t, ok)
	})

	t.Run("ReturnsNotOKIfNonStringKeyDoesntExist", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, string]()
		_, ok := m.Get(123)
		assert.False(t, ok)
	})

	t.Run("ReturnsOKIfKeyExists", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		_, ok := m.Get("foo")
		assert.True(t, ok)
	})

	t.Run("ReturnsValueForKey", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		value, _ := m.Get("foo")
		assert.Equal(t, "bar", value)
	})

	t.Run("ReturnsDynamicValueForKey", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "baz")
		value, _ := m.Get("foo")
		assert.Equal(t, "baz", value)
	})

	t.Run("KeyDoesntExistOnNonEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "baz")
		_, ok := m.Get("bar")
		assert.False(t, ok)
	})

	t.Run("ValueForKeyDoesntExistOnNonEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "baz")
		value, _ := m.Get("bar")
		assert.Empty(t, value)
	})
}

func TestSet(t *testing.T) {
	t.Run("ReturnsTrueIfStringKeyIsNew", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		ok := m.Set("foo", "bar")
		assert.True(t, ok)
	})

	t.Run("ReturnsTrueIfNonStringKeyIsNew", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, string]()
		ok := m.Set(123, "bar")
		assert.True(t, ok)
	})

	t.Run("ValueCanBeNonString", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		ok := m.Set(123, true)
		assert.True(t, ok)
	})

	t.Run("ReturnsFalseIfKeyIsNotNew", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		ok := m.Set("foo", "bar")
		assert.False(t, ok)
	})

	t.Run("SetThreeDifferentKeys", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		m.Set("baz", "qux")
		ok := m.Set("quux", "corge")
		assert.True(t, ok)
	})
}

func TestLen(t *testing.T) {
	t.Run("EmptyMapIsZeroLen", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		assert.Equal(t, 0, m.Len())
	})

	t.Run("SingleElementIsLenOne", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		m.Set(123, true)
		assert.Equal(t, 1, m.Len())
	})

	t.Run("ThreeElements", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		m.Set(1, true)
		m.Set(2, true)
		m.Set(3, true)
		assert.Equal(t, 3, m.Len())
	})
}

func TestKeys(t *testing.T) {
	t.Run("EmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		assert.Empty(t, m.Keys())
	})

	t.Run("OneElement", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		m.Set(1, true)
		assert.Equal(t, []int{1}, m.Keys())
	})

	t.Run("RetainsOrder", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		for i := 1; i < 10; i++ {
			m.Set(i, true)
		}
		assert.Equal(t,
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			m.Keys())
	})

	t.Run("ReplacingKeyDoesntChangeOrder", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, bool]()
		m.Set("foo", true)
		m.Set("bar", true)
		m.Set("foo", false)
		assert.Equal(t,
			[]string{"foo", "bar"},
			m.Keys())
	})

	t.Run("KeysAfterDelete", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, bool]()
		m.Set("foo", true)
		m.Set("bar", true)
		m.Delete("foo")
		assert.Equal(t, []string{"bar"}, m.Keys())
	})
}

func TestDelete(t *testing.T) {
	t.Run("KeyDoesntExistReturnsFalse", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		assert.False(t, m.Delete("foo"))
	})

	t.Run("KeyDoesExist", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, any]()
		m.Set("foo", nil)
		assert.True(t, m.Delete("foo"))
	})

	t.Run("KeyNoLongerExists", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, any]()
		m.Set("foo", nil)
		m.Delete("foo")
		_, exists := m.Get("foo")
		assert.False(t, exists)
	})

	t.Run("KeyDeleteIsIsolated", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, any]()
		m.Set("foo", nil)
		m.Set("bar", nil)
		m.Delete("foo")
		_, exists := m.Get("bar")
		assert.True(t, exists)
	})
}

func TestOrderedMap_Front(t *testing.T) {
	t.Run("NilOnEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		assert.Nil(t, m.Front())
	})

	t.Run("NilOnEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		m.Set(1, true)
		assert.NotNil(t, m.Front())
	})
}

func TestOrderedMap_Back(t *testing.T) {
	t.Run("NilOnEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		assert.Nil(t, m.Back())
	})

	t.Run("NilOnEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		m.Set(1, true)
		assert.NotNil(t, m.Back())
	})
}

func TestOrderedMap_Copy(t *testing.T) {
	t.Run("ReturnsEqualButNotSame", func(t *testing.T) {
		key, value := 1, "a value"
		m := orderedmap.NewOrderedMap[int, string]()
		m.Set(key, value)

		m2 := m.Copy()
		m2.Set(key, "a different value")

		assert.Equal(t, m.Len(), m2.Len(), "not all elements are copied")
		assert.Equal(t, value, m.GetElement(key).Value)
	})
}

func TestGetElement(t *testing.T) {
	t.Run("ReturnsElementForKey", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "bar")

		var results []any
		element := m.GetElement("foo")
		if element != nil {
			results = append(results, element.Key, element.Value)
		}

		assert.Equal(t, []any{"foo", "bar"}, results)
	})

	t.Run("ElementForKeyDoesntExistOnNonEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "baz")
		element := m.GetElement("bar")
		assert.Nil(t, element)
	})
}
