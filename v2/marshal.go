package orderedmap

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func (m *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, 0, m.Len()<<3))
	b.WriteByte('{')
	for el := m.Front(); el != nil; el = el.Next() {
		keyValue := reflect.ValueOf(el.Key)
		keyStr, err := resolveKeyName(keyValue)
		if err != nil {
			return nil, fmt.Errorf("error resolving key: %w", err)
		}
		key, err := json.Marshal(keyStr)
		if err != nil {
			return nil, fmt.Errorf("error marshaling key %v: %w", el.Key, err)
		}
		b.Write(key)
		b.WriteByte(':')
		value, err := json.Marshal(el.Value)
		if err != nil {
			return nil, fmt.Errorf("error marshaling value for key %v: %w", el.Key, err)
		}
		b.Write(value)
		if el.Next() != nil {
			b.WriteByte(',')
		}
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

func resolveKeyName(k reflect.Value) (string, error) {
	if k.Kind() == reflect.String {
		return k.String(), nil
	}
	if tm, ok := k.Interface().(encoding.TextMarshaler); ok {
		if k.Kind() == reflect.Pointer && k.IsNil() {
			return "", nil
		}
		buf, err := tm.MarshalText()
		return string(buf), err
	}
	switch k.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(k.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(k.Uint(), 10), nil
	}
	return "", &json.UnsupportedTypeError{Type: k.Type()}
}
