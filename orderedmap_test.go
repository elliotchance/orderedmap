package orderedmap_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/elliotchance/orderedmap"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderedMap(t *testing.T) {
	m := orderedmap.NewOrderedMap()
	assert.IsType(t, &orderedmap.OrderedMap{}, m)
}

func TestGet(t *testing.T) {
	t.Run("ReturnsNotOKIfStringKeyDoesntExist", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		_, ok := m.Get("foo")
		assert.False(t, ok)
	})

	t.Run("ReturnsNotOKIfNonStringKeyDoesntExist", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		_, ok := m.Get(123)
		assert.False(t, ok)
	})

	t.Run("ReturnsOKIfKeyExists", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", "bar")
		_, ok := m.Get("foo")
		assert.True(t, ok)
	})

	t.Run("ReturnsValueForKey", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", "bar")
		value, _ := m.Get("foo")
		assert.Equal(t, "bar", value)
	})

	t.Run("ReturnsDynamicValueForKey", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", "baz")
		value, _ := m.Get("foo")
		assert.Equal(t, "baz", value)
	})

	t.Run("KeyDoesntExistOnNonEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", "baz")
		_, ok := m.Get("bar")
		assert.False(t, ok)
	})

	t.Run("ValueForKeyDoesntExistOnNonEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", "baz")
		value, _ := m.Get("bar")
		assert.Nil(t, value)
	})

	t.Run("Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkOrderedMap_Get(1))
		res4 := testing.Benchmark(benchmarkOrderedMap_Get(4))

		// O(1) would mean that res4 should take about the same time as res1,
		// because we are accessing the same amount of elements, just on
		// different sized maps.

		assert.InDelta(t,
			res1.NsPerOp(), res4.NsPerOp(),
			0.5*float64(res1.NsPerOp()))
	})
}

func TestSet(t *testing.T) {
	t.Run("ReturnsTrueIfStringKeyIsNew", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		ok := m.Set("foo", "bar")
		assert.True(t, ok)
	})

	t.Run("ReturnsTrueIfNonStringKeyIsNew", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		ok := m.Set(123, "bar")
		assert.True(t, ok)
	})

	t.Run("ValueCanBeNonString", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		ok := m.Set(123, true)
		assert.True(t, ok)
	})

	t.Run("ReturnsFalseIfKeyIsNotNew", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", "bar")
		ok := m.Set("foo", "bar")
		assert.False(t, ok)
	})

	t.Run("SetThreeDifferentKeys", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", "bar")
		m.Set("baz", "qux")
		ok := m.Set("quux", "corge")
		assert.True(t, ok)
	})

	t.Run("Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkOrderedMap_Set(1))
		res4 := testing.Benchmark(benchmarkOrderedMap_Set(4))

		// O(1) would mean that res4 should take about 4 times longer than res1
		// because we are doing 4 times the amount of Set operations. Allow for
		// a wide margin, but not too wide that it would permit the inflection
		// to O(n^2).

		assert.InDelta(t,
			4*res1.NsPerOp(), res4.NsPerOp(),
			2*float64(res1.NsPerOp()))
	})
}

func TestLen(t *testing.T) {
	t.Run("EmptyMapIsZeroLen", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		assert.Equal(t, 0, m.Len())
	})

	t.Run("SingleElementIsLenOne", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set(123, true)
		assert.Equal(t, 1, m.Len())
	})

	t.Run("ThreeElements", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set(1, true)
		m.Set(2, true)
		m.Set(3, true)
		assert.Equal(t, 3, m.Len())
	})

	t.Run("Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkOrderedMap_Len(1))
		res4 := testing.Benchmark(benchmarkOrderedMap_Len(4))

		// O(1) would mean that res4 should take about the same time as res1,
		// because we are accessing the same amount of elements, just on
		// different sized maps.

		assert.InDelta(t,
			res1.NsPerOp(), res4.NsPerOp(),
			0.5*float64(res1.NsPerOp()))
	})
}

func TestKeys(t *testing.T) {
	t.Run("EmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		assert.Empty(t, m.Keys())
	})

	t.Run("OneElement", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set(1, true)
		assert.Equal(t, []interface{}{1}, m.Keys())
	})

	t.Run("RetainsOrder", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		for i := 1; i < 10; i++ {
			m.Set(i, true)
		}
		assert.Equal(t,
			[]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9},
			m.Keys())
	})

	t.Run("ReplacingKeyDoesntChangeOrder", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", true)
		m.Set("bar", true)
		m.Set("foo", false)
		assert.Equal(t,
			[]interface{}{"foo", "bar"},
			m.Keys())
	})

	t.Run("KeysAfterDelete", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", true)
		m.Set("bar", true)
		m.Delete("foo")
		assert.Equal(t, []interface{}{"bar"}, m.Keys())
	})

	t.Run("Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkOrderedMap_Keys(1))
		res4 := testing.Benchmark(benchmarkOrderedMap_Keys(4))

		// O(1) would mean that res4 should take about 4 times longer than res1
		// because we are doing 4 times the amount of Set/Delete operations.
		// Allow for a wide margin, but not too wide that it would permit the
		// inflection to O(n^2).

		assert.InDelta(t,
			4*res1.NsPerOp(), res4.NsPerOp(),
			float64(res4.NsPerOp()))
	})
}

func TestDelete(t *testing.T) {
	t.Run("KeyDoesntExistReturnsFalse", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		assert.False(t, m.Delete("foo"))
	})

	t.Run("KeyDoesExist", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", nil)
		assert.True(t, m.Delete("foo"))
	})

	t.Run("KeyNoLongerExists", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", nil)
		m.Delete("foo")
		_, exists := m.Get("foo")
		assert.False(t, exists)
	})

	t.Run("KeyDeleteIsIsolated", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", nil)
		m.Set("bar", nil)
		m.Delete("foo")
		_, exists := m.Get("bar")
		assert.True(t, exists)
	})

	t.Run("Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkOrderedMap_Delete(1))
		res4 := testing.Benchmark(benchmarkOrderedMap_Delete(4))

		// O(1) would mean that res4 should take about 4 times longer than res1
		// because we are doing 4 times the amount of Set/Delete operations.
		// Allow for a wide margin, but not too wide that it would permit the
		// inflection to O(n^2).

		assert.InDelta(t,
			4*res1.NsPerOp(), res4.NsPerOp(),
			float64(res4.NsPerOp()))
	})
}

func TestOrderedMap_Front(t *testing.T) {
	t.Run("NilOnEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		assert.Nil(t, m.Front())
	})

	t.Run("NilOnEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set(1, true)
		assert.NotNil(t, m.Front())
	})
}

func TestOrderedMap_Back(t *testing.T) {
	t.Run("NilOnEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		assert.Nil(t, m.Back())
	})

	t.Run("NilOnEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set(1, true)
		assert.NotNil(t, m.Back())
	})
}

func TestGetElement(t *testing.T) {
	t.Run("ReturnsNotOKIfStringKeyDoesntExist", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		_, ok := m.GetElement("foo")
		assert.False(t, ok)
	})

	t.Run("ReturnsNotOKIfNonStringKeyDoesntExist", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		_, ok := m.GetElement(123)
		assert.False(t, ok)
	})

	t.Run("ReturnsOKIfKeyExists", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", "bar")
		_, ok := m.GetElement("foo")
		assert.True(t, ok)
	})

	t.Run("ReturnsElementForKey", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", "bar")

		var results []interface{}
		element, _ := m.GetElement("foo")
		if element != nil {
			results = append(results, element.Key, element.Value)
		}

		assert.Equal(t, []interface{}{"foo", "bar"}, results)
	})

	t.Run("ReturnsDynamicValueForKey", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", "bar")

		var results []interface{}
		element, _ := m.GetElement("foo")
		if element != nil {
			results = append(results, element.Key, element.Value)
		}

		assert.Equal(t, []interface{}{"foo", "bar"}, results)
	})

	t.Run("KeyDoesntExistOnNonEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", "baz")
		_, ok := m.GetElement("bar")
		assert.False(t, ok)
	})

	t.Run("ElementForKeyDoesntExistOnNonEmptyMap", func(t *testing.T) {
		m := orderedmap.NewOrderedMap()
		m.Set("foo", "baz")
		element, _ := m.GetElement("bar")
		assert.Nil(t, element)
	})

	t.Run("Performance", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		res1 := testing.Benchmark(benchmarkOrderedMap_GetElement(1))
		res4 := testing.Benchmark(benchmarkOrderedMap_GetElement(4))

		// O(1) would mean that res4 should take about the same time as res1,
		// because we are accessing the same amount of elements, just on
		// different sized maps.

		assert.InDelta(t,
			res1.NsPerOp(), res4.NsPerOp(),
			0.5*float64(res1.NsPerOp()))
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
		m := orderedmap.NewOrderedMap()
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
	m := orderedmap.NewOrderedMap()
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

func benchmarkMap_GetElement(multiplier int) func(b *testing.B) {
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

func BenchmarkMap_GetElement(b *testing.B) {
	benchmarkMap_GetElement(1)(b)
}

func benchmarkOrderedMap_GetElement(multiplier int) func(b *testing.B) {
	m := orderedmap.NewOrderedMap()
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
	m := orderedmap.NewOrderedMap()
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
		m := orderedmap.NewOrderedMap()
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
	m := orderedmap.NewOrderedMap()
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, key := range m.Keys() {
				_, v := m.Get(key)
				nothing(v)
			}
		}
	}
}

func BenchmarkOrderedMap_Iterate(b *testing.B) {
	benchmarkOrderedMap_Iterate(1)(b)
}

func benchmarkOrderedMap_Keys(multiplier int) func(b *testing.B) {
	m := orderedmap.NewOrderedMap()
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Keys()
		}
	}
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
		m := orderedmap.NewOrderedMap()
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
	m := orderedmap.NewOrderedMap()
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

func benchmarkMapString_GetElement(multiplier int) func(b *testing.B) {
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

func BenchmarkMapString_GetElement(b *testing.B) {
	benchmarkMapString_GetElement(1)(b)
}

func benchmarkOrderedMapString_GetElement(multiplier int) func(b *testing.B) {
	m := orderedmap.NewOrderedMap()
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
		m := orderedmap.NewOrderedMap()
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
	m := orderedmap.NewOrderedMap()
	a := "12345678"
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, key := range m.Keys() {
				_, v := m.Get(key)
				nothing(v)
			}
		}
	}
}

func BenchmarkOrderedMapString_Iterate(b *testing.B) {
	benchmarkOrderedMapString_Iterate(1)(b)
}

func BenchmarkOrderedMap_Keys(b *testing.B) {
	benchmarkOrderedMap_Keys(1)(b)
}

func ExampleNewOrderedMap() {
	m := orderedmap.NewOrderedMap()

	m.Set("foo", "bar")
	m.Set("qux", 1.23)
	m.Set(123, true)

	m.Delete("qux")

	for _, key := range m.Keys() {
		value, _ := m.Get(key)
		fmt.Println(key, value)
	}
}

func ExampleOrderedMap_Front() {
	m := orderedmap.NewOrderedMap()
	m.Set(1, true)
	m.Set(2, true)

	for el := m.Front(); el != nil; el = el.Next() {
		fmt.Println(el)
	}
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
			m := orderedmap.NewOrderedMap()
			for i := 0; i < 10000000; i++ {
				m.Set(i, true)
			}
		}
	}
}

func BenchmarkBigOrderedMap_Set(b *testing.B) {
	benchmarkBigOrderedMap_Set()(b)
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
	m := orderedmap.NewOrderedMap()
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

func benchmarkBigMap_GetElement() func(b *testing.B) {
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

func BenchmarkBigMap_GetElement(b *testing.B) {
	benchmarkBigMap_GetElement()(b)
}

func benchmarkBigOrderedMap_GetElement() func(b *testing.B) {
	m := orderedmap.NewOrderedMap()
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
	m := orderedmap.NewOrderedMap()
	for i := 0; i < 10000000; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, key := range m.Keys() {
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

func BenchmarkBigMapString_Set(b *testing.B) {
	benchmarkBigMapString_Set()(b)
}

func benchmarkBigOrderedMapString_Set() func(b *testing.B) {
	return func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			m := orderedmap.NewOrderedMap()
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
	m := orderedmap.NewOrderedMap()
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

func benchmarkBigMapString_GetElement() func(b *testing.B) {
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

func BenchmarkBigMapString_GetElement(b *testing.B) {
	benchmarkBigMapString_GetElement()(b)
}

func benchmarkBigOrderedMapString_GetElement() func(b *testing.B) {
	m := orderedmap.NewOrderedMap()
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
	m := orderedmap.NewOrderedMap()
	a := "12345678"
	for i := 0; i < 10000000; i++ {
		m.Set(a+strconv.Itoa(i), true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, key := range m.Keys() {
				_, v := m.Get(key)
				nothing(v)
			}
		}
	}
}

func BenchmarkBigOrderedMapString_Iterate(b *testing.B) {
	benchmarkBigOrderedMapString_Iterate()(b)
}

func BenchmarkAll(b *testing.B) {
	b.Run("BenchmarkOrderedMap_Keys", BenchmarkOrderedMap_Keys)

	b.Run("BenchmarkOrderedMap_Set", BenchmarkOrderedMap_Set)
	b.Run("BenchmarkMap_Set", BenchmarkMap_Set)
	b.Run("BenchmarkOrderedMap_Get", BenchmarkOrderedMap_Get)
	b.Run("BenchmarkMap_Get", BenchmarkMap_Get)
	b.Run("BenchmarkOrderedMap_GetElement", BenchmarkOrderedMap_GetElement)
	b.Run("BenchmarkMap_GetElement", BenchmarkMap_GetElement)
	b.Run("BenchmarkOrderedMap_Delete", BenchmarkOrderedMap_Delete)
	b.Run("BenchmarkMap_Delete", BenchmarkMap_Delete)
	b.Run("BenchmarkOrderedMap_Iterate", BenchmarkOrderedMap_Iterate)
	b.Run("BenchmarkMap_Iterate", BenchmarkMap_Iterate)

	b.Run("BenchmarkBigMap_Set", BenchmarkBigMap_Set)
	b.Run("BenchmarkBigOrderedMap_Set", BenchmarkBigOrderedMap_Set)
	b.Run("BenchmarkBigMap_Get", BenchmarkBigMap_Get)
	b.Run("BenchmarkBigOrderedMap_Get", BenchmarkBigOrderedMap_Get)
	b.Run("BenchmarkBigMap_GetElement", BenchmarkBigMap_GetElement)
	b.Run("BenchmarkBigOrderedMap_GetElement", BenchmarkBigOrderedMap_GetElement)
	b.Run("BenchmarkBigOrderedMap_Iterate", BenchmarkBigOrderedMap_Iterate)
	b.Run("BenchmarkBigMap_Iterate", BenchmarkBigMap_Iterate)

	b.Run("BenchmarkOrderedMapString_Set", BenchmarkOrderedMapString_Set)
	b.Run("BenchmarkMapString_Set", BenchmarkMapString_Set)
	b.Run("BenchmarkOrderedMapString_Get", BenchmarkOrderedMapString_Get)
	b.Run("BenchmarkMapString_Get", BenchmarkMapString_Get)
	b.Run("BenchmarkOrderedMapString_GetElement", BenchmarkOrderedMapString_GetElement)
	b.Run("BenchmarkMapString_GetElement", BenchmarkMapString_GetElement)
	b.Run("BenchmarkOrderedMapString_Delete", BenchmarkOrderedMapString_Delete)
	b.Run("BenchmarkMapString_Delete", BenchmarkMapString_Delete)
	b.Run("BenchmarkOrderedMapString_Iterate", BenchmarkOrderedMapString_Iterate)
	b.Run("BenchmarkMapString_Iterate", BenchmarkMapString_Iterate)

	b.Run("BenchmarkBigMapString_Set", BenchmarkBigMapString_Set)
	b.Run("BenchmarkBigOrderedMapString_Set", BenchmarkBigOrderedMapString_Set)
	b.Run("BenchmarkBigMapString_Get", BenchmarkBigMapString_Get)
	b.Run("BenchmarkBigOrderedMapString_Get", BenchmarkBigOrderedMapString_Get)
	b.Run("BenchmarkBigMapString_GetElement", BenchmarkBigMapString_GetElement)
	b.Run("BenchmarkBigOrderedMapString_GetElement", BenchmarkBigOrderedMapString_GetElement)
	b.Run("BenchmarkBigOrderedMapString_Iterate", BenchmarkBigOrderedMapString_Iterate)
	b.Run("BenchmarkBigMapString_Iterate", BenchmarkBigMapString_Iterate)
}
