package yarf

import (
	"reflect"
)

//NewParam creates a new key value param from input
func NewParam(key string, value interface{}) Param {
	return Param{key, value}
}

//Param is a key/value entry and a struct which implements helper methods to help with retrial of data types from value.
type Param struct {
	key   string
	value interface{}
}

// Key returns the key of the key/value pair
func (m *Param) Key() string {
	return m.key
}

// Value returns the value of the key/value pair
func (m *Param) Value() interface{} {
	return m.value
}

// String returns value as a string, if possible
func (m *Param) String() (string, bool) {
	if m.value == nil {
		return "", false
	}

	str, ok := m.value.(string)
	return str, ok
}

// StringOr returns value as a string, otherwise the provided default
func (m *Param) StringOr(defaultTo string) string {
	str, ok := m.String()
	if ok {
		return str
	}
	return defaultTo
}

// StringArr returns value as a []string, if possible
func (m *Param) StringArr() ([]string, bool) {
	if m.value == nil {
		return nil, false
	}

	var res []string
	res, ok := m.value.([]string)
	if ok {
		return res, ok
	}

	arr, ok := m.value.([]interface{})

	if !ok {
		return res, ok
	}
	ok = true
	res = make([]string, len(arr))
	for i, val := range arr {
		var o bool
		res[i], o = val.(string)
		if !o {
			ok = false
		}
	}
	return res, ok
}

// StringArrOr returns value as a []string, otherwise the provided default
func (m *Param) StringArrOr(defaultTo []string) []string {
	arr, ok := m.StringArr()

	if ok {
		return arr
	}
	return defaultTo
}

// Int returns value as a int64, if possible
func (m *Param) Int() (int64, bool) {
	return toInt(m.value)
}

func toInt(num interface{}) (int64, bool) {

	if num == nil {
		return 0, false
	}

	var i int64
	ok := false

	switch num.(type) {
	case int, int8, int16, int32, int64:
		a := reflect.ValueOf(num).Int() // a has type int64
		return a, true
	case uint, uint8, uint16, uint32, uint64:
		a := reflect.ValueOf(num).Uint() // a has type uint64
		return int64(a), true
	case float64:
		f, ok := num.(float64)
		return int64(f), ok
	case float32:
		f, ok := num.(float32)
		return int64(f), ok
	}

	return i, ok

}

// IntOr returns value as a int64, otherwise the provided default
func (m *Param) IntOr(def int64) int64 {
	i, ok := m.Int()
	if ok {
		return i
	}
	return def
}

// IntArr returns value as a []int64, if possible
func (m *Param) IntArr() ([]int64, bool) {
	if m.value == nil {
		return nil, false
	}
	arr, ok := m.value.([]interface{})
	var res []int64
	if !ok {
		return res, ok
	}
	res = make([]int64, len(arr))
	ok = true
	for i, val := range arr {
		var o bool
		res[i], o = val.(int64)
		if !o {
			ok = false
		}
	}
	return res, ok
}

// IntArrOr returns value as a []int64, otherwise the provided default
func (m *Param) IntArrOr(def []int64) []int64 {
	arr, ok := m.IntArr()

	if ok {
		return arr
	}
	return def
}

// Float returns value as a float64, if possible
func (m *Param) Float() (float64, bool) {
	if m.value == nil {
		return 0.0, false
	}

	i, ok := m.value.(float64)
	return i, ok
}

// FloatOr returns value as a float64, otherwise the provided default
func (m *Param) FloatOr(def float64) float64 {
	i, ok := m.Float()
	if ok {
		return i
	}
	return def
}

// FloatArr returns value as a []float64, if possible
func (m *Param) FloatArr() ([]float64, bool) {
	if m.value == nil {
		return nil, false
	}
	arr, ok := m.value.([]interface{})
	var res []float64
	if !ok {
		return res, ok
	}
	res = make([]float64, len(arr))
	ok = true
	for i, val := range arr {
		var o bool
		res[i], o = val.(float64)
		if !o {
			ok = false
		}
	}
	return res, ok
}

// FloatArrOr returns value as a []float64, otherwise the provided default
func (m *Param) FloatArrOr(def []float64) []float64 {
	arr, ok := m.FloatArr()

	if ok {
		return arr
	}
	return def
}

// Bool returns value as a bool, if possible
func (m *Param) Bool() (bool, bool) {
	if m.value == nil {
		return false, false
	}

	i, ok := m.value.(bool)
	return i, ok
}

// BoolOr returns value as a bool, otherwise the provided default
func (m *Param) BoolOr(def bool) bool {
	i, ok := m.Bool()
	if ok {
		return i
	}
	return def
}

// BoolArr returns value as a []bool, if possible
func (m *Param) BoolArr() ([]bool, bool) {
	if m.value == nil {
		return nil, false
	}
	arr, ok := m.value.([]interface{})
	var res []bool
	if !ok {
		return res, ok
	}
	res = make([]bool, len(arr))
	ok = true
	for i, val := range arr {
		var o bool
		res[i], o = val.(bool)
		if !o {
			ok = false
		}
	}
	return res, ok
}

// BoolArrOr returns value as a []bool, otherwise the provided default
func (m *Param) BoolArrOr(def []bool) []bool {
	arr, ok := m.BoolArr()

	if ok {
		return arr
	}
	return def
}
