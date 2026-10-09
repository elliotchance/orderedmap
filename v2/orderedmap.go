package orderedmap

import (
	"bytes"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type OrderedMap[K comparable, V any] struct {
	kv map[K]*Element[K, V]
	ll list[K, V]
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		kv: make(map[K]*Element[K, V]),
	}
}

// NewOrderedMapWithCapacity creates a map with enough pre-allocated space to
// hold the specified number of elements.
func NewOrderedMapWithCapacity[K comparable, V any](capacity int) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		kv: make(map[K]*Element[K, V], capacity),
	}
}

func NewOrderedMapWithElements[K comparable, V any](els ...*Element[K, V]) *OrderedMap[K, V] {
	om := NewOrderedMapWithCapacity[K, V](len(els))
	for _, el := range els {
		om.Set(el.Key, el.Value)
	}
	return om
}

// Get returns the value for a key. If the key does not exist, the second return
// parameter will be false and the value will be nil.
func (m *OrderedMap[K, V]) Get(key K) (value V, ok bool) {
	v, ok := m.kv[key]
	if ok {
		value = v.Value
	}

	return
}

// Set will set (or replace) a value for a key. If the key was new, then true
// will be returned. The returned value will be false if the value was replaced
// (even if the value was the same).
func (m *OrderedMap[K, V]) Set(key K, value V) bool {
	if m.kv == nil {
		m.kv = make(map[K]*Element[K, V])
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

// ReplaceKey replaces an existing key with a new key while preserving order of
// the value. This function will return true if the operation was successful, or
// false if 'originalKey' is not found OR 'newKey' already exists (which would be an overwrite).
func (m *OrderedMap[K, V]) ReplaceKey(originalKey, newKey K) bool {
	element, originalExists := m.kv[originalKey]
	_, newKeyExists := m.kv[newKey]
	if originalExists && !newKeyExists {
		delete(m.kv, originalKey)
		m.kv[newKey] = element
		element.Key = newKey
		return true
	}
	return false
}

// GetOrDefault returns the value for a key. If the key does not exist, returns
// the default value instead.
func (m *OrderedMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	if value, ok := m.kv[key]; ok {
		return value.Value
	}

	return defaultValue
}

// GetElement returns the element for a key. If the key does not exist, the
// pointer will be nil.
func (m *OrderedMap[K, V]) GetElement(key K) *Element[K, V] {
	element, ok := m.kv[key]
	if ok {
		return element
	}

	return nil
}

// Len returns the number of elements in the map.
func (m *OrderedMap[K, V]) Len() int {
	return len(m.kv)
}

// Keys returns all of the keys in the order they were inserted. If a key was
// replaced it will retain the same position. To ensure most recently set keys
// are always at the end you must always Delete before Set.
func (m *OrderedMap[K, V]) Keys() (keys []K) {
	keys = make([]K, 0, m.Len())
	for el := m.Front(); el != nil; el = el.Next() {
		keys = append(keys, el.Key)
	}
	return keys
}

// Delete will remove a key from the map. It will return true if the key was
// removed (the key did exist).
func (m *OrderedMap[K, V]) Delete(key K) (didDelete bool) {
	element, ok := m.kv[key]
	if ok {
		m.ll.Remove(element)
		delete(m.kv, key)
	}

	return ok
}

// Front will return the element that is the first (oldest Set element). If
// there are no elements this will return nil.
func (m *OrderedMap[K, V]) Front() *Element[K, V] {
	return m.ll.Front()
}

// Back will return the element that is the last (most recent Set element). If
// there are no elements this will return nil.
func (m *OrderedMap[K, V]) Back() *Element[K, V] {
	return m.ll.Back()
}

// Copy returns a new OrderedMap with the same elements.
// Using Copy while there are concurrent writes may mangle the result.
func (m *OrderedMap[K, V]) Copy() *OrderedMap[K, V] {
	m2 := NewOrderedMapWithCapacity[K, V](m.Len())
	for el := m.Front(); el != nil; el = el.Next() {
		m2.Set(el.Key, el.Value)
	}
	return m2
}

// Has checks if a key exists in the map.
func (m *OrderedMap[K, V]) Has(key K) bool {
	_, exists := m.kv[key]
	return exists
}

// MarshalJSON implements the json.Marshaler interface.
func (m *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
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
func (m *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("orderedmap: UnmarshalJSON on nil pointer")
	}
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return errors.New("unexpected end of JSON input")
	}
	if bytes.Equal(trimmed, []byte("null")) {
		m.kv = nil
		m.ll = list[K, V]{}
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
		m.kv = make(map[K]*Element[K, V])
	}

	isInterfaceVal := reflect.TypeOf((*V)(nil)).Elem().Kind() == reflect.Interface

	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return err
		}
		keyStr, ok := keyTok.(string)
		if !ok {
			return fmt.Errorf("expected string key, got %v", keyTok)
		}

		key, err := unmarshalKey[K](keyStr)
		if err != nil {
			return err
		}

		var val V
		if isInterfaceVal {
			decoded, err := decodeValue(dec)
			if err != nil {
				return err
			}
			val = decoded.(V)
		} else {
			if err := dec.Decode(&val); err != nil {
				return err
			}
		}

		m.Set(key, val)
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

func marshalKey[K comparable](key K) ([]byte, error) {
	if tm, ok := any(key).(encoding.TextMarshaler); ok {
		text, err := tm.MarshalText()
		if err != nil {
			return nil, err
		}
		return json.Marshal(string(text))
	}
	if s, ok := any(key).(string); ok {
		return json.Marshal(s)
	}
	return json.Marshal(fmt.Sprint(key))
}

func unmarshalKey[K comparable](keyStr string) (K, error) {
	var key K
	if _, ok := any(key).(string); ok {
		return any(keyStr).(K), nil
	}
	keyJSON, err := json.Marshal(keyStr)
	if err != nil {
		return key, err
	}
	dummyJSON := append(append([]byte("{"), keyJSON...), []byte(":null}")...)
	var dummy map[K]interface{}
	if err := json.Unmarshal(dummyJSON, &dummy); err != nil {
		return key, err
	}
	for k := range dummy {
		key = k
	}
	return key, nil
}

func decodeValue(dec *json.Decoder) (any, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	switch t := tok.(type) {
	case json.Delim:
		switch t {
		case '{':
			child := NewOrderedMap[string, any]()
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
			var list []any
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

