# 🔃 github.com/elliotchance/orderedmap [![GoDoc](https://godoc.org/github.com/elliotchance/orderedmap?status.svg)](https://godoc.org/github.com/elliotchance/orderedmap) [![Build Status](https://travis-ci.org/elliotchance/orderedmap.svg?branch=master)](https://travis-ci.org/elliotchance/orderedmap)

## Installation

```bash
go get -u github.com/elliotchance/orderedmap
```

## Basic Usage

An `*OrderedMap` is a high performance ordered map that maintains amortized O(1)
for `Set`, `Get`, `Delete` and `Len`:

```go
m := orderedmap.NewOrderedMap()

m.Set("foo", "bar")
m.Set("qux", 1.23)
m.Set(123, true)

m.Delete("qux")
```

Internally an `*OrderedMap` uses a combination of a map and linked list.

## Iterating

Be careful using `Keys()` as it will create a copy of all of the keys so it's
only suitable for a small number of items:

```go
for _, key := range m.Keys() {
	value, _:= m.Get(key)
	fmt.Println(key, value)
}
```

For larger maps you should use `Front()` or `Back()` to iterate per element:

```go
// Iterate through all elements from oldest to newest:
for el := m.Front(); el != nil; el = el.Next() {
    fmt.Println(el.Key, el.Value)
}

// You can also use Back and Prev to iterate in reverse:
for el := m.Back(); el != nil; el = el.Prev() {
    fmt.Println(el.Key, el.Value)
}
```

The iterator is safe to use bidirectionally, and will return `nil` once it goes
beyond the first or last item.

If the map is changing while the iteration is in-flight it may produce
unexpected behavior.
