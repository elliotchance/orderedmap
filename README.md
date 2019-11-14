# 🔃 github.com/elliotchance/orderedmap [![GoDoc](https://godoc.org/github.com/elliotchance/orderedmap?status.svg)](https://godoc.org/github.com/elliotchance/orderedmap)

The `orderedmap` package provides a high performance ordered map in Go:

```go
m := orderedmap.NewOrderedMap()

m.Set("foo", "bar")
m.Set("qux", 1.23)
m.Set(123, true)

m.Delete("qux")

for _, key := range m.Keys() {
	value, _:= m.Get(key)
	fmt.Println(key, value)
}
```

Internally an `*OrderedMap` uses a combination of a map and linked list to
maintain amortized O(1) for `Set`, `Get`, `Delete` and `Len`. 

See the full documentation at
[https://godoc.org/github.com/elliotchance/orderedmap](https://godoc.org/github.com/elliotchance/orderedmap).
