package orderedmap_test

import (
	"encoding/json"
	"slices"
	"strconv"
	"testing"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestReplaceKey(t *testing.T) {
	t.Run("ReturnsFalseIfOriginalKeyDoesntExist", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		assert.False(t, m.ReplaceKey("foo", "bar"))
	})

	t.Run("ReturnsFalseIfNewKeyAlreadyExists", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		m.Set("baz", "qux")
		assert.False(t, m.ReplaceKey("foo", "baz"))
		assert.Equal(t, []string{"foo", "baz"}, slices.Collect(m.Keys()))
	})

	t.Run("ReturnsTrueIfOnlyOriginalKeyExists", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		assert.True(t, m.ReplaceKey("foo", "baz"))

		// Now validate the "replacement" was a success.
		el := m.GetElement("baz")
		require.NotNil(t, el)
		assert.Equal(t, "bar", el.Value)
		assert.Equal(t, "baz", el.Key)

		v, ok := m.Get("baz")
		assert.True(t, ok)
		assert.Equal(t, "bar", v)
		assert.Equal(t, []string{"baz"}, slices.Collect(m.Keys()))
		assert.Equal(t, 1, m.Len())

		_, ok = m.Get("foo") // original key
		assert.False(t, ok)
	})

	t.Run("KeyMaintainsOrderWhenReplaced", func(t *testing.T) {
		count := 100
		// Build a larger map to help validate that the order is not coincidental.
		m := orderedmap.NewOrderedMap[int, int]()
		for i := 0; i < count; i++ {
			m.Set(i, i)
		}
		// Rename the middle 50-60 elements to 100+ current
		for i := 50; i < 60; i++ {
			assert.True(t, m.ReplaceKey(i, i+100))
		}

		// ensure length is maintained
		assert.Equal(t, count, m.Len())

		// Validate the order is maintained.
		for i, key := range slices.Collect(m.Keys()) {
			if i >= 50 && i < 60 {
				assert.Equal(t, i+100, key)
			} else {
				assert.Equal(t, i, key)
			}
		}
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
		assert.Empty(t, slices.Collect(m.Keys()))
	})

	t.Run("OneElement", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		m.Set(1, true)
		assert.Equal(t, []int{1}, slices.Collect(m.Keys()))
	})

	t.Run("RetainsOrder", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		for i := 1; i < 10; i++ {
			m.Set(i, true)
		}
		assert.Equal(t,
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			slices.Collect(m.Keys()))
	})

	t.Run("ReplacingKeyDoesntChangeOrder", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, bool]()
		m.Set("foo", true)
		m.Set("bar", true)
		m.Set("foo", false)
		assert.Equal(t,
			[]string{"foo", "bar"},
			slices.Collect(m.Keys()))
	})

	t.Run("KeysAfterDelete", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, bool]()
		m.Set("foo", true)
		m.Set("bar", true)
		m.Delete("foo")
		assert.Equal(t, []string{"bar"}, slices.Collect(m.Keys()))
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

func TestSetAndGet(t *testing.T) {
	t.Run("FourBoolElements", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		expected := map[int]bool{1: true, 3: false, 5: false, 4: true}
		for k, v := range expected {
			m.Set(k, v)
		}
		for k, v := range expected {
			w, ok := m.Get(k)
			assert.True(t, ok)
			assert.Equal(t, v, w)
		}
	})
}

func TestIterations(t *testing.T) {
	type Element struct {
		Key   int
		Value bool
	}
	t.Run("FourBoolElements", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, bool]()
		expected := []Element{{5, true}, {3, false}, {1, false}, {4, true}}
		for _, v := range expected {
			m.Set(v.Key, v.Value)
		}
		element := m.Front()
		for i := 0; i < len(expected); i++ {
			assert.NotNil(t, element)
			assert.Equal(t, expected[i].Key, element.Key)
			assert.Equal(t, expected[i].Value, element.Value)
			element = element.Next()
		}
		assert.Nil(t, element)
	})
}

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
		for key, value := range m.AllFromFront() {
			assert.Equal(t, expected[i].Key, key)
			assert.Equal(t, expected[i].Value, value)
			i++
		}
	})

	t.Run("ReverseIterator", func(t *testing.T) {
		i := len(expected) - 1
		for key, value := range m.AllFromBack() {
			assert.Equal(t, expected[i].Key, key)
			assert.Equal(t, expected[i].Value, value)
			i--
		}
	})
}

func TestOrderedMap_Has(t *testing.T) {
	t.Run("ReturnsFalseIfKeyDoesNotExist", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		assert.False(t, m.Has("foo"))
	})

	t.Run("ReturnsTrueIfKeyExists", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		assert.True(t, m.Has("foo"))
	})

	t.Run("KeyDoesNotExistAfterDelete", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		m.Delete("foo")
		assert.False(t, m.Has("foo"))
	})
}

type textKey struct {
	val string
}

func (k textKey) MarshalText() ([]byte, error) {
	return []byte("text_" + k.val), nil
}

func (k *textKey) UnmarshalText(text []byte) error {
	k.val = string(text)
	return nil
}

func TestOrderedMap_MarshalJSON(t *testing.T) {
	t.Run("NilMap", func(t *testing.T) {
		var m *orderedmap.OrderedMap[string, int]
		b, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, "null", string(b))
	})

	t.Run("EmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, int]()
		b, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, "{}", string(b))
	})

	t.Run("SingleElement", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		b, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, `{"foo":"bar"}`, string(b))
	})

	t.Run("PreservesInsertionOrder", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, int]()
		m.Set("zebra", 1)
		m.Set("apple", 2)
		m.Set("mango", 3)
		b, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, `{"zebra":1,"apple":2,"mango":3}`, string(b))
	})

	t.Run("ReplacingKeyPreservesPosition", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("a", 10)
		b, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, `{"a":10,"b":2}`, string(b))
	})

	t.Run("IntegerKeys", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[int, string]()
		m.Set(3, "three")
		m.Set(1, "one")
		m.Set(2, "two")
		b, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, `{"3":"three","1":"one","2":"two"}`, string(b))
	})

	t.Run("TextMarshalerKeys", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[textKey, string]()
		m.Set(textKey{val: "b"}, "val_b")
		m.Set(textKey{val: "a"}, "val_a")
		b, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, `{"text_b":"val_b","text_a":"val_a"}`, string(b))
	})

	t.Run("NestedOrderedMap", func(t *testing.T) {
		inner := orderedmap.NewOrderedMap[string, int]()
		inner.Set("z", 100)
		inner.Set("a", 200)

		outer := orderedmap.NewOrderedMap[string, any]()
		outer.Set("id", "test")
		outer.Set("nested", inner)

		b, err := json.Marshal(outer)
		assert.NoError(t, err)
		assert.Equal(t, `{"id":"test","nested":{"z":100,"a":200}}`, string(b))
	})

	t.Run("MarshalIndent", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, int]()
		m.Set("b", 2)
		m.Set("a", 1)
		b, err := json.MarshalIndent(m, "", "  ")
		assert.NoError(t, err)
		expected := "{\n  \"b\": 2,\n  \"a\": 1\n}"
		assert.Equal(t, expected, string(b))
	})
}

func TestOrderedMap_UnmarshalJSON(t *testing.T) {
	t.Run("EmptyInputError", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, int]()
		err := json.Unmarshal([]byte(""), m)
		assert.Error(t, err)
	})

	t.Run("NilPointerError", func(t *testing.T) {
		var m *orderedmap.OrderedMap[string, int]
		err := m.UnmarshalJSON([]byte(`{"a":1}`))
		assert.Error(t, err)
	})

	t.Run("NullInput", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, string]()
		m.Set("foo", "bar")
		err := json.Unmarshal([]byte("null"), m)
		assert.NoError(t, err)
		assert.Equal(t, 0, m.Len())
	})

	t.Run("EmptyObject", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, int]()
		err := json.Unmarshal([]byte("{}"), m)
		assert.NoError(t, err)
		assert.Equal(t, 0, m.Len())
	})

	t.Run("PreservesOrderStringKeys", func(t *testing.T) {
		input := `{"z":1,"a":2,"m":3}`
		m := orderedmap.NewOrderedMap[string, int]()
		err := json.Unmarshal([]byte(input), m)
		assert.NoError(t, err)
		assert.Equal(t, []string{"z", "a", "m"}, slices.Collect(m.Keys()))
		v, ok := m.Get("a")
		assert.True(t, ok)
		assert.Equal(t, 2, v)
	})

	t.Run("PreservesOrderIntegerKeys", func(t *testing.T) {
		input := `{"30":"thirty","10":"ten","20":"twenty"}`
		m := orderedmap.NewOrderedMap[int, string]()
		err := json.Unmarshal([]byte(input), m)
		assert.NoError(t, err)
		assert.Equal(t, []int{30, 10, 20}, slices.Collect(m.Keys()))
		v, ok := m.Get(10)
		assert.True(t, ok)
		assert.Equal(t, "ten", v)
	})

	t.Run("DuplicateKeysPreservesOriginalOrderAndUpdatesValue", func(t *testing.T) {
		input := `{"a":1,"b":2,"a":3}`
		m := orderedmap.NewOrderedMap[string, int]()
		err := json.Unmarshal([]byte(input), m)
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b"}, slices.Collect(m.Keys()))
		v, ok := m.Get("a")
		assert.True(t, ok)
		assert.Equal(t, 3, v)
	})

	t.Run("NestedObjectsWithAny", func(t *testing.T) {
		input := `{"outer":{"z":10,"y":20,"x":30}}`
		m := orderedmap.NewOrderedMap[string, any]()
		err := json.Unmarshal([]byte(input), m)
		assert.NoError(t, err)

		outerVal, ok := m.Get("outer")
		assert.True(t, ok)
		childMap, isOM := outerVal.(*orderedmap.OrderedMap[string, any])
		assert.True(t, isOM)
		assert.Equal(t, []string{"z", "y", "x"}, slices.Collect(childMap.Keys()))
		valY, _ := childMap.Get("y")
		assert.Equal(t, float64(20), valY)
	})

	t.Run("ArrayWithNestedObjects", func(t *testing.T) {
		input := `{"list":[{"k2":"v2","k1":"v1"}]}`
		m := orderedmap.NewOrderedMap[string, any]()
		err := json.Unmarshal([]byte(input), m)
		assert.NoError(t, err)

		listVal, ok := m.Get("list")
		assert.True(t, ok)
		list := listVal.([]any)
		assert.Len(t, list, 1)

		elem1 := list[0].(*orderedmap.OrderedMap[string, any])
		assert.Equal(t, []string{"k2", "k1"}, slices.Collect(elem1.Keys()))
	})

	t.Run("InvalidSyntax", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, int]()
		err := json.Unmarshal([]byte(`{"a":}`), m)
		assert.Error(t, err)
	})

	t.Run("InvalidTypeNotObject", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, int]()
		err := json.Unmarshal([]byte(`[1, 2, 3]`), m)
		assert.Error(t, err)
	})

	t.Run("TrailingData", func(t *testing.T) {
		m := orderedmap.NewOrderedMap[string, int]()
		err := json.Unmarshal([]byte(`{"a": 1} trailing`), m)
		assert.Error(t, err)
	})
}

func TestOrderedMap_JSONRoundTrip(t *testing.T) {
	t.Run("ComplexNestedRoundTrip", func(t *testing.T) {
		originalJSON := `{"name":"test","config":{"timeout":30,"retries":3,"tags":["fast","ci"]},"active":true}`
		m := orderedmap.NewOrderedMap[string, any]()
		err := json.Unmarshal([]byte(originalJSON), m)
		assert.NoError(t, err)

		marshaled, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, originalJSON, string(marshaled))
	})

	t.Run("StructContainingOrderedMap", func(t *testing.T) {
		type Envelope struct {
			Title   string                               `json:"title"`
			Payload *orderedmap.OrderedMap[string, int] `json:"payload"`
		}

		payload := orderedmap.NewOrderedMap[string, int]()
		payload.Set("z", 1)
		payload.Set("a", 2)
		env := Envelope{
			Title:   "envelope",
			Payload: payload,
		}

		b, err := json.Marshal(env)
		assert.NoError(t, err)
		assert.Equal(t, `{"title":"envelope","payload":{"z":1,"a":2}}`, string(b))

		var envDecoded Envelope
		err = json.Unmarshal(b, &envDecoded)
		assert.NoError(t, err)
		assert.Equal(t, "envelope", envDecoded.Title)
		assert.Equal(t, []string{"z", "a"}, slices.Collect(envDecoded.Payload.Keys()))
	})
}

func benchmarkMap_Set(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := make(map[int]bool)
		for i := 0; i < b.N*multiplier; i++ {
			m[i] = true
		}
	}
}

func BenchmarkMap_Set(b *testing.B) {
	benchmarkMap_Set(1)(b)
}

func benchmarkOrderedMap_Set(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := orderedmap.NewOrderedMap[int, bool]()
		for i := 0; i < b.N*multiplier; i++ {
			m.Set(i, true)
		}
	}
}

func BenchmarkOrderedMap_Set(b *testing.B) {
	benchmarkOrderedMap_Set(1)(b)
}

func benchmarkMap_Get(multiplier int) func(b *testing.B) {
	m := make(map[int]bool)
	for i := 0; i < 1000*multiplier; i++ {
		m[i] = true
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m[i%1000*multiplier]
		}
	}
}

func BenchmarkMap_Get(b *testing.B) {
	benchmarkMap_Get(1)(b)
}

func benchmarkOrderedMap_Get(multiplier int) func(b *testing.B) {
	m := orderedmap.NewOrderedMap[int, bool]()
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Get(i % 1000 * multiplier)
		}
	}
}

func BenchmarkOrderedMap_Get(b *testing.B) {
	benchmarkOrderedMap_Get(1)(b)
}

func benchmarkOrderedMap_GetElement(multiplier int) func(b *testing.B) {
	m := orderedmap.NewOrderedMap[int, bool]()
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.GetElement(i % 1000 * multiplier)
		}
	}
}

func BenchmarkOrderedMap_GetElement(b *testing.B) {
	benchmarkOrderedMap_GetElement(1)(b)
}

var tempInt int

func benchmarkOrderedMap_Len(multiplier int) func(b *testing.B) {
	m := orderedmap.NewOrderedMap[int, bool]()
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		var temp int
		for i := 0; i < b.N; i++ {
			temp = m.Len()
		}

		// prevent compiler from optimising Len away.
		tempInt = temp
	}
}

func BenchmarkOrderedMap_Len(b *testing.B) {
	benchmarkOrderedMap_Len(1)(b)
}

func benchmarkMap_Delete(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := make(map[int]bool)
		for i := 0; i < b.N*multiplier; i++ {
			m[i] = true
		}

		for i := 0; i < b.N; i++ {
			delete(m, i)
		}
	}
}

func BenchmarkMap_Delete(b *testing.B) {
	benchmarkMap_Delete(1)(b)
}

func benchmarkOrderedMap_Delete(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := orderedmap.NewOrderedMap[int, bool]()
		for i := 0; i < b.N*multiplier; i++ {
			m.Set(i, true)
		}

		for i := 0; i < b.N; i++ {
			m.Delete(i)
		}
	}
}

func BenchmarkOrderedMap_Delete(b *testing.B) {
	benchmarkOrderedMap_Delete(1)(b)
}

func benchmarkMap_Iterate(multiplier int) func(b *testing.B) {
	m := make(map[int]bool)
	for i := 0; i < 1000*multiplier; i++ {
		m[i] = true
	}
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range m {
				nothing(v)
			}
		}
	}
}

func BenchmarkMap_Iterate(b *testing.B) {
	benchmarkMap_Iterate(1)(b)
}

func benchmarkOrderedMap_Iterate(multiplier int) func(b *testing.B) {
	m := orderedmap.NewOrderedMap[int, bool]()
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, key := range slices.Collect(m.Keys()) {
				_, v := m.Get(key)
				nothing(v)
			}
		}
	}
}

func BenchmarkOrderedMap_Iterate(b *testing.B) {
	benchmarkOrderedMap_Iterate(1)(b)
}

func benchmarkMap_Has(multiplier int) func(b *testing.B) {
	m := make(map[int]bool)
	for i := 0; i < 1000*multiplier; i++ {
		m[i] = true
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m[i%1000*multiplier]
		}
	}
}

func BenchmarkMap_Has(b *testing.B) {
	benchmarkMap_Has(1)(b)
}

func benchmarkOrderedMap_Has(multiplier int) func(b *testing.B) {
	m := orderedmap.NewOrderedMap[int, bool]()
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Has(i % 1000 * multiplier)
		}
	}
}

func BenchmarkOrderedMap_Has(b *testing.B) {
	benchmarkOrderedMap_Has(1)(b)
}

func benchmarkMapString_Set(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := make(map[string]bool)
		a := "12345678"
		for i := 0; i < b.N*multiplier; i++ {
			m[a+strconv.Itoa(i)] = true
		}
	}
}

func BenchmarkMapString_Set(b *testing.B) {
	benchmarkMapString_Set(1)(b)
}

func benchmarkOrderedMapString_Set(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := orderedmap.NewOrderedMap[string, bool]()
		a := "12345678"
		for i := 0; i < b.N*multiplier; i++ {
			m.Set(a+strconv.Itoa(i), true)
		}
	}
}

func BenchmarkOrderedMapString_Set(b *testing.B) {
	benchmarkOrderedMapString_Set(1)(b)
}

func benchmarkMapString_Get(multiplier int) func(b *testing.B) {
	m := make(map[string]bool)
	a := "12345678"
	for i := 0; i < 1000*multiplier; i++ {
		m[a+strconv.Itoa(i)] = true
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m[a+strconv.Itoa(i%1000*multiplier)]
		}
	}
}

func BenchmarkMapString_Get(b *testing.B) {
	benchmarkMapString_Get(1)(b)
}

func benchmarkOrderedMapString_Get(multiplier int) func(b *testing.B) {
	m := orderedmap.NewOrderedMap[string, bool]()
	a := "12345678"
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Get(a + strconv.Itoa(i%1000*multiplier))
		}
	}
}

func BenchmarkOrderedMapString_Get(b *testing.B) {
	benchmarkOrderedMapString_Get(1)(b)
}

func benchmarkOrderedMapString_GetElement(multiplier int) func(b *testing.B) {
	m := orderedmap.NewOrderedMap[string, bool]()
	a := "12345678"
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.GetElement(a + strconv.Itoa(i%1000*multiplier))
		}
	}
}

func BenchmarkOrderedMapString_GetElement(b *testing.B) {
	benchmarkOrderedMapString_GetElement(1)(b)
}

func benchmarkMapString_Delete(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := make(map[string]bool)
		a := "12345678"
		for i := 0; i < b.N*multiplier; i++ {
			m[a+strconv.Itoa(i)] = true
		}

		for i := 0; i < b.N; i++ {
			delete(m, a+strconv.Itoa(i))
		}
	}
}

func BenchmarkMapString_Delete(b *testing.B) {
	benchmarkMapString_Delete(1)(b)
}

func benchmarkOrderedMapString_Delete(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := orderedmap.NewOrderedMap[string, bool]()
		a := "12345678"
		for i := 0; i < b.N*multiplier; i++ {
			m.Set(a+strconv.Itoa(i), true)
		}

		for i := 0; i < b.N; i++ {
			m.Delete(a + strconv.Itoa(i))
		}
	}
}

func BenchmarkOrderedMapString_Delete(b *testing.B) {
	benchmarkOrderedMapString_Delete(1)(b)
}

func benchmarkMapString_Iterate(multiplier int) func(b *testing.B) {
	m := make(map[string]bool)
	a := "12345678"
	for i := 0; i < 1000*multiplier; i++ {
		m[a+strconv.Itoa(i)] = true
	}
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range m {
				nothing(v)
			}
		}
	}
}

func BenchmarkMapString_Iterate(b *testing.B) {
	benchmarkMapString_Iterate(1)(b)
}

func benchmarkOrderedMapString_Iterate(multiplier int) func(b *testing.B) {
	m := orderedmap.NewOrderedMap[string, bool]()
	a := "12345678"
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, key := range slices.Collect(m.Keys()) {
				_, v := m.Get(key)
				nothing(v)
			}
		}
	}
}

func BenchmarkOrderedMapString_Iterate(b *testing.B) {
	benchmarkOrderedMapString_Iterate(1)(b)
}

func benchmarkMapString_Has(multiplier int) func(b *testing.B) {
	m := make(map[string]bool)
	a := "12345678"
	for i := 0; i < 1000*multiplier; i++ {
		m[a+strconv.Itoa(i)] = true
	}
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m[a+strconv.Itoa(i%1000*multiplier)]
		}
	}
}

func BenchmarkMapString_Has(b *testing.B) {
	benchmarkMapString_Has(1)(b)
}

func benchmarkOrderedMapString_Has(multiplier int) func(b *testing.B) {
	m := orderedmap.NewOrderedMap[string, bool]()
	a := "12345678"
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Has(a + strconv.Itoa(i%1000*multiplier))
		}
	}
}

func BenchmarkOrderedMapString_Has(b *testing.B) {
	benchmarkOrderedMapString_Has(1)(b)
}

func nothing(v interface{}) {
	v = false
}

func benchmarkBigMap_Set() func(b *testing.B) {
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			m := make(map[int]bool)
			for i := 0; i < 10000000; i++ {
				m[i] = true
			}
		}
	}
}

func BenchmarkBigMap_Set(b *testing.B) {
	benchmarkBigMap_Set()(b)
}

func benchmarkBigOrderedMap_Set() func(b *testing.B) {
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			m := orderedmap.NewOrderedMap[int, bool]()
			for i := 0; i < 10000000; i++ {
				m.Set(i, true)
			}
		}
	}
}

func BenchmarkBigOrderedMap_Set(b *testing.B) {
	benchmarkBigOrderedMap_Set()(b)
}

func benchmarkBigMapWithCapacity_Set() func(b *testing.B) {
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			m := orderedmap.NewOrderedMapWithCapacity[int, bool](10000000)
			for i := 0; i < 10000000; i++ {
				m.Set(i, true)
			}
		}
	}
}

func BenchmarkBigMapWithCapacity_Set(b *testing.B) {
	benchmarkBigMapWithCapacity_Set()(b)
}

func benchmarkBigMap_Get() func(b *testing.B) {
	m := make(map[int]bool)
	for i := 0; i < 10000000; i++ {
		m[i] = true
	}

	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				_ = m[i]
			}
		}
	}
}

func BenchmarkBigMap_Get(b *testing.B) {
	benchmarkBigMap_Get()(b)
}

func benchmarkBigOrderedMap_Get() func(b *testing.B) {
	m := orderedmap.NewOrderedMap[int, bool]()
	for i := 0; i < 10000000; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				m.Get(i)
			}
		}
	}
}

func BenchmarkBigOrderedMap_Get(b *testing.B) {
	benchmarkBigOrderedMap_Get()(b)
}

func benchmarkBigOrderedMap_GetElement() func(b *testing.B) {
	m := orderedmap.NewOrderedMap[int, bool]()
	for i := 0; i < 10000000; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				m.GetElement(i)
			}
		}
	}
}

func BenchmarkBigOrderedMap_GetElement(b *testing.B) {
	benchmarkBigOrderedMap_GetElement()(b)
}

func benchmarkBigMap_Iterate() func(b *testing.B) {
	m := make(map[int]bool)
	for i := 0; i < 10000000; i++ {
		m[i] = true
	}
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range m {
				nothing(v)
			}
		}
	}
}

func BenchmarkBigMap_Iterate(b *testing.B) {
	benchmarkBigMap_Iterate()(b)
}

func benchmarkBigOrderedMap_Iterate() func(b *testing.B) {
	m := orderedmap.NewOrderedMap[int, bool]()
	for i := 0; i < 10000000; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, key := range slices.Collect(m.Keys()) {
				_, v := m.Get(key)
				nothing(v)
			}
		}
	}
}

func BenchmarkBigOrderedMap_Iterate(b *testing.B) {
	benchmarkBigOrderedMap_Iterate()(b)
}

func benchmarkBigMapString_Set() func(b *testing.B) {
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			m := make(map[string]bool)
			a := "1234567"
			for i := 0; i < 10000000; i++ {
				m[a+strconv.Itoa(i)] = true
			}
		}
	}
}

func benchmarkBigMap_Has() func(b *testing.B) {
	m := make(map[int]bool)
	for i := 0; i < 10000000; i++ {
		m[i] = true
	}
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				_ = m[i]
			}
		}
	}
}

func BenchmarkBigMap_Has(b *testing.B) {
	benchmarkBigMap_Has()(b)
}

func benchmarkBigOrderedMap_Has() func(b *testing.B) {
	m := orderedmap.NewOrderedMap[int, bool]()
	for i := 0; i < 10000000; i++ {
		m.Set(i, true)
	}
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				m.Has(i)
			}
		}
	}
}

func BenchmarkBigOrderedMap_Has(b *testing.B) {
	benchmarkBigOrderedMap_Has()(b)
}

func BenchmarkBigMapString_Set(b *testing.B) {
	benchmarkBigMapString_Set()(b)
}

func benchmarkBigOrderedMapString_Set() func(b *testing.B) {
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			m := orderedmap.NewOrderedMap[string, bool]()
			a := "1234567"
			for i := 0; i < 10000000; i++ {
				m.Set(a+strconv.Itoa(i), true)
			}
		}
	}
}

func BenchmarkBigOrderedMapString_Set(b *testing.B) {
	benchmarkBigOrderedMapString_Set()(b)
}

func benchmarkBigMapString_Get() func(b *testing.B) {
	m := make(map[string]bool)
	a := "1234567"
	for i := 0; i < 10000000; i++ {
		m[a+strconv.Itoa(i)] = true
	}

	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				_ = m[a+strconv.Itoa(i)]
			}
		}
	}
}

func BenchmarkBigMapString_Get(b *testing.B) {
	benchmarkBigMapString_Get()(b)
}

func benchmarkBigOrderedMapString_Get() func(b *testing.B) {
	m := orderedmap.NewOrderedMap[string, bool]()
	a := "1234567"
	for i := 0; i < 10000000; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}

	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				m.Get(a + strconv.Itoa(i))
			}
		}
	}
}

func BenchmarkBigOrderedMapString_Get(b *testing.B) {
	benchmarkBigOrderedMapString_Get()(b)
}

func benchmarkBigOrderedMapString_GetElement() func(b *testing.B) {
	m := orderedmap.NewOrderedMap[string, bool]()
	a := "1234567"
	for i := 0; i < 10000000; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}

	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				m.GetElement(a + strconv.Itoa(i))
			}
		}
	}
}

func BenchmarkBigOrderedMapString_GetElement(b *testing.B) {
	benchmarkBigOrderedMapString_GetElement()(b)
}

func benchmarkBigMapString_Iterate() func(b *testing.B) {
	m := make(map[string]bool)
	a := "12345678"
	for i := 0; i < 10000000; i++ {
		m[a+strconv.Itoa(i)] = true
	}
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range m {
				nothing(v)
			}
		}
	}
}

func BenchmarkBigMapString_Iterate(b *testing.B) {
	benchmarkBigMapString_Iterate()(b)
}

func benchmarkBigOrderedMapString_Iterate() func(b *testing.B) {
	m := orderedmap.NewOrderedMap[string, bool]()
	a := "12345678"
	for i := 0; i < 10000000; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, key := range slices.Collect(m.Keys()) {
				_, v := m.Get(key)
				nothing(v)
			}
		}
	}
}

func BenchmarkBigOrderedMapString_Iterate(b *testing.B) {
	benchmarkBigOrderedMapString_Iterate()(b)
}

func benchmarkBigMapString_Has() func(b *testing.B) {
	m := make(map[string]bool)
	a := "12345678"
	for i := 0; i < 10000000; i++ {
		m[a+strconv.Itoa(i)] = true
	}
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				_ = m[a+strconv.Itoa(i)]
			}
		}
	}
}

func BenchmarkBigMapString_Has(b *testing.B) {
	benchmarkBigMapString_Has()(b)
}

func benchmarkBigOrderedMapString_Has() func(b *testing.B) {
	m := orderedmap.NewOrderedMap[string, bool]()
	a := "12345678"
	for i := 0; i < 10000000; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < 10000000; i++ {
				m.Has(a + strconv.Itoa(i))
			}
		}
	}
}

func BenchmarkBigOrderedMapString_Has(b *testing.B) {
	benchmarkBigOrderedMapString_Has()(b)
}

func BenchmarkAll(b *testing.B) {
	b.Run("BenchmarkOrderedMap_Set", BenchmarkOrderedMap_Set)
	b.Run("BenchmarkMap_Set", BenchmarkMap_Set)
	b.Run("BenchmarkOrderedMap_Get", BenchmarkOrderedMap_Get)
	b.Run("BenchmarkMap_Get", BenchmarkMap_Get)
	b.Run("BenchmarkOrderedMap_GetElement", BenchmarkOrderedMap_GetElement)
	b.Run("BenchmarkOrderedMap_Delete", BenchmarkOrderedMap_Delete)
	b.Run("BenchmarkMap_Delete", BenchmarkMap_Delete)
	b.Run("BenchmarkOrderedMap_Iterate", BenchmarkOrderedMap_Iterate)
	b.Run("BenchmarkMap_Iterate", BenchmarkMap_Iterate)
	b.Run("BenchmarkOrderedMap_Has", BenchmarkOrderedMap_Has)
	b.Run("BenchmarkMap_Has", BenchmarkMap_Has)

	b.Run("BenchmarkBigMap_Set", BenchmarkBigMap_Set)
	b.Run("BenchmarkBigOrderedMap_Set", BenchmarkBigOrderedMap_Set)
	b.Run("BenchmarkBigMap_Get", BenchmarkBigMap_Get)
	b.Run("BenchmarkBigOrderedMap_Get", BenchmarkBigOrderedMap_Get)
	b.Run("BenchmarkBigOrderedMap_GetElement",
		BenchmarkBigOrderedMap_GetElement)
	b.Run("BenchmarkBigOrderedMap_Iterate", BenchmarkBigOrderedMap_Iterate)
	b.Run("BenchmarkBigMap_Iterate", BenchmarkBigMap_Iterate)
	b.Run("BenchmarkBigMap_Has", BenchmarkBigMap_Has)
	b.Run("BenchmarkBigOrderedMap_Has", BenchmarkBigOrderedMap_Has)

	b.Run("BenchmarkOrderedMapString_Set", BenchmarkOrderedMapString_Set)
	b.Run("BenchmarkMapString_Set", BenchmarkMapString_Set)
	b.Run("BenchmarkOrderedMapString_Get", BenchmarkOrderedMapString_Get)
	b.Run("BenchmarkMapString_Get", BenchmarkMapString_Get)
	b.Run("BenchmarkOrderedMapString_GetElement",
		BenchmarkOrderedMapString_GetElement)
	b.Run("BenchmarkOrderedMapString_Delete", BenchmarkOrderedMapString_Delete)
	b.Run("BenchmarkMapString_Delete", BenchmarkMapString_Delete)
	b.Run("BenchmarkOrderedMapString_Iterate",
		BenchmarkOrderedMapString_Iterate)
	b.Run("BenchmarkMapString_Iterate", BenchmarkMapString_Iterate)
	b.Run("BenchmarkMapString_Has", BenchmarkMapString_Has)
	b.Run("BenchmarkOrderedMapString_Has", BenchmarkOrderedMapString_Has)

	b.Run("BenchmarkBigMapString_Set", BenchmarkBigMapString_Set)
	b.Run("BenchmarkBigOrderedMapString_Set", BenchmarkBigOrderedMapString_Set)
	b.Run("BenchmarkBigMapString_Get", BenchmarkBigMapString_Get)
	b.Run("BenchmarkBigOrderedMapString_Get", BenchmarkBigOrderedMapString_Get)
	b.Run("BenchmarkBigOrderedMapString_GetElement",
		BenchmarkBigOrderedMapString_GetElement)
	b.Run("BenchmarkBigOrderedMapString_Iterate",
		BenchmarkBigOrderedMapString_Iterate)
	b.Run("BenchmarkBigMapString_Iterate", BenchmarkBigMapString_Iterate)
	b.Run("BenchmarkBigMapString_Has", BenchmarkBigMapString_Has)
	b.Run("BenchmarkBigOrderedMapString_Has", BenchmarkBigOrderedMapString_Has)
}
