package orderedmap_test

import (
	"fmt"
	"github.com/elliotchance/orderedmap"
	"github.com/stretchr/testify/assert"
	"testing"
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

func benchmarkOrderedMap_Get(multiplier int) func(b *testing.B) {
	m := orderedmap.NewOrderedMap()
	for i := 0; i < 1000*multiplier; i++ {
		m.Set(i, true)
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Get(1000 * multiplier % b.N)
		}
	}
}

func BenchmarkOrderedMap_Get(b *testing.B) {
	benchmarkOrderedMap_Get(1)(b)
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

func benchmarkOrderedMap_Delete(multiplier int) func(b *testing.B) {
	return func(b *testing.B) {
		m := orderedmap.NewOrderedMap()
		for i := 0; i < b.N*multiplier; i++ {
			m.Set(i, true)
		}

		for i := 0; i < b.N; i++ {
			m.Delete(b.N)
		}
	}
}

func BenchmarkOrderedMap_Delete(b *testing.B) {
	benchmarkOrderedMap_Delete(1)(b)
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
