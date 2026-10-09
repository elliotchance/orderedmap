package orderedmap

import (
	"bytes"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type OrderedMap struct {
	kv map[interface{}]*Element
	ll list
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		kv: make(map[interface{}]*Element),
	}
}

// NewOrderedMapWithCapacity creates a map with enough pre-allocated space to
// hold the specified number of elements.
func NewOrderedMapWithCapacity(capacity int) *OrderedMap {
	return &OrderedMap{
		kv: make(map[interface{}]*Element, capacity),
	}
}

// Get returns the value for a key. If the key does not exist, the second return
// parameter will be false and the value will be nil.
func (m *OrderedMap) Get(key interface{}) (interface{}, bool) {
	element, ok := m.kv[key]
	if ok {
		return element.Value, true
	}

	return nil, false
}

// Set will set (or replace) a value for a key. If the key was new, then true
// will be returned. The returned value will be false if the value was replaced
// (even if the value was the same).
func (m *OrderedMap) Set(key, value interface{}) bool {
	if m.kv == nil {
		m.kv = make(map[interface{}]*Element)
	}
	_, alreadyExist := m.kv[key]
	if alreadyExist {
		m.kv[key].Value = value
		return false
	}

	element := m.ll.PushBack(key, value)
	m.kv[key] = element
	return true
}

// GetOrDefault returns the value for a key. If the key does not exist, returns
// the default value instead.
func (m *OrderedMap) GetOrDefault(key, defaultValue interface{}) interface{} {
	if element, ok := m.kv[key]; ok {
		return element.Value
	}

	return defaultValue
}

// GetElement returns the element for a key. If the key does not exist, the
// pointer will be nil.
func (m *OrderedMap) GetElement(key interface{}) *Element {
	element, ok := m.kv[key]
	if ok {
		return element
	}

	return nil
}

// Len returns the number of elements in the map.
func (m *OrderedMap) Len() int {
	return len(m.kv)
}

// Keys returns all of the keys in the order they were inserted. If a key was
// replaced it will retain the same position. To ensure most recently set keys
// are always at the end you must always Delete before Set.
func (m *OrderedMap) Keys() (keys []interface{}) {
	keys = make([]interface{}, 0, m.Len())
	for el := m.Front(); el != nil; el = el.Next() {
		keys = append(keys, el.Key)
	}
	return keys
}

// Delete will remove a key from the map. It will return true if the key was
// removed (the key did exist).
func (m *OrderedMap) Delete(key interface{}) (didDelete bool) {
	element, ok := m.kv[key]
	if ok {
		m.ll.Remove(element)
		delete(m.kv, key)
	}

	return ok
}

// Front will return the element that is the first (oldest Set element). If
// there are no elements this will return nil.
func (m *OrderedMap) Front() *Element {
	return m.ll.Front()
}

// Back will return the element that is the last (most recent Set element). If
// there are no elements this will return nil.
func (m *OrderedMap) Back() *Element {
	return m.ll.Back()
}

// Copy returns a new OrderedMap with the same elements.
// Using Copy while there are concurrent writes may mangle the result.
func (m *OrderedMap) Copy() *OrderedMap {
	m2 := NewOrderedMapWithCapacity(m.Len())

	for el := m.Front(); el != nil; el = el.Next() {
		m2.Set(el.Key, el.Value)
	}

	return m2
}

// Has checks if a key exists in the map.
func (m *OrderedMap) Has(key interface{}) bool {
	_, exists := m.kv[key]
	return exists
}

// MarshalJSON implements the json.Marshaler interface.
func (m *OrderedMap) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	if m.Len() == 0 {
		return []byte("{}"), nil
	}

	var buf bytes.Buffer
	buf.WriteByte('{')
	first := true
	for el := m.Front(); el != nil; el = el.Next() {
		if !first {
			buf.WriteByte(',')
		}
		first = false

		keyBytes, err := marshalKey(el.Key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteByte(':')

		valBytes, err := json.Marshal(el.Value)
		if err != nil {
			return nil, err
		}
		buf.Write(valBytes)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *OrderedMap) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("orderedmap: UnmarshalJSON on nil pointer")
	}
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return errors.New("unexpected end of JSON input")
	}
	if bytes.Equal(trimmed, []byte("null")) {
		m.kv = nil
		m.ll = list{}
		return nil
	}

	dec := json.NewDecoder(bytes.NewReader(trimmed))
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '{' {
		return fmt.Errorf("expected '{', got %v", tok)
	}

	if m.kv == nil {
		m.kv = make(map[interface{}]*Element)
	}

	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return err
		}
		keyStr, ok := keyTok.(string)
		if !ok {
			return fmt.Errorf("expected string key, got %v", keyTok)
		}

		val, err := decodeValue(dec)
		if err != nil {
			return err
		}
		m.Set(keyStr, val)
	}

	tok, err = dec.Token()
	if err != nil {
		return err
	}
	delim, ok = tok.(json.Delim)
	if !ok || delim != '}' {
		return fmt.Errorf("expected '}', got %v", tok)
	}

	_, err = dec.Token()
	if err != io.EOF {
		return fmt.Errorf("unexpected data after JSON object")
	}

	return nil
}

func marshalKey(key interface{}) ([]byte, error) {
	if key == nil {
		return json.Marshal("null")
	}
	if tm, ok := key.(encoding.TextMarshaler); ok {
		text, err := tm.MarshalText()
		if err != nil {
			return nil, err
		}
		return json.Marshal(string(text))
	}
	if s, ok := key.(string); ok {
		return json.Marshal(s)
	}
	return json.Marshal(fmt.Sprint(key))
}

func decodeValue(dec *json.Decoder) (interface{}, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	switch t := tok.(type) {
	case json.Delim:
		switch t {
		case '{':
			child := NewOrderedMap()
			for dec.More() {
				kt, err := dec.Token()
				if err != nil {
					return nil, err
				}
				keyStr, ok := kt.(string)
				if !ok {
					return nil, fmt.Errorf("expected string key, got %v", kt)
				}
				val, err := decodeValue(dec)
				if err != nil {
					return nil, err
				}
				child.Set(keyStr, val)
			}
			closeTok, err := dec.Token()
			if err != nil {
				return nil, err
			}
			if d, ok := closeTok.(json.Delim); !ok || d != '}' {
				return nil, fmt.Errorf("expected '}', got %v", closeTok)
			}
			return child, nil
		case '[':
			var list []interface{}
			for dec.More() {
				elem, err := decodeValue(dec)
				if err != nil {
					return nil, err
				}
				list = append(list, elem)
			}
			closeTok, err := dec.Token()
			if err != nil {
				return nil, err
			}
			if d, ok := closeTok.(json.Delim); !ok || d != ']' {
				return nil, fmt.Errorf("expected ']', got %v", closeTok)
			}
			return list, nil
		default:
			return nil, fmt.Errorf("unexpected delimiter: %v", t)
		}
	default:
		return t, nil
	}
}

