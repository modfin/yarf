package yarf

import (
	"reflect"
)

var uintType = reflect.TypeOf(uint64(0))
var intType = reflect.TypeOf(int64(0))
var floatType = reflect.TypeOf(float64(0))
var stringType = reflect.TypeOf(string(""))
var boolType = reflect.TypeOf(false)

type converter func(in interface{}) (interface{}, bool)

func untypedFloat(in interface{}) (interface{}, bool) {
	return toFloat(in)
}
func untypedUint(in interface{}) (interface{}, bool) {
	return toUint(in)
}
func untypedInt(in interface{}) (interface{}, bool) {
	return toInt(in)
}
func untypedBool(in interface{}) (res interface{}, ok bool) {
	return toBool(in)
}
func untypedString(in interface{}) (res interface{}, ok bool) {
	return toString(in)
}

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

// IsNil returns true if value is nil
func (m *Param) IsNil() bool {
	return m.value == nil
}

// IsSlice returns true if value is a array
func (m *Param) IsSlice() bool {

	if m.value == nil {
		return false
	}
	return reflect.TypeOf(m.value).Kind() == reflect.Slice
}

// String returns value as a string, if possible
func (m *Param) String() (string, bool) {
	return toString(m.value)
}

// StringOr returns value as a string, otherwise the provided default
func (m *Param) StringOr(defaultTo string) string {
	str, ok := m.String()
	if ok {
		return str
	}
	return defaultTo
}

// StringSlice returns value as a []string, if possible
func (m *Param) StringSlice() ([]string, bool) {
	if m.value == nil {
		return nil, false
	}

	var res []string
	res, ok := m.value.([]string)
	if ok {
		return res, true
	}

	r, ok := toSliceOf(m.value, stringType, untypedString)
	if !ok {
		return nil, false
	}
	res, ok = r.([]string)
	if ok {
		return res, true
	}
	return nil, false
}

// StringSliceOr returns value as a []string, otherwise the provided default
func (m *Param) StringSliceOr(defaultTo []string) []string {
	arr, ok := m.StringSlice()

	if ok {
		return arr
	}
	return defaultTo
}

// Uint returns value as a uint64, if possible
func (m *Param) Uint() (uint64, bool) {
	return toUint(m.value)
}

// UintOr returns value as a uint64, otherwise the provided default
func (m *Param) UintOr(def uint64) uint64 {
	i, ok := m.Uint()
	if ok {
		return i
	}
	return def
}

// UintSlice returns value as a []uint64, if possible
func (m *Param) UintSlice() ([]uint64, bool) {
	if m.value == nil {
		return nil, false
	}
	var res []uint64
	res, ok := m.value.([]uint64)

	if ok {
		return res, true
	}

	r, ok := toSliceOf(m.value, uintType, untypedUint)
	if !ok {
		return nil, false
	}
	res, ok = r.([]uint64)
	if ok {
		return res, true
	}
	return nil, false
}

// UintSliceOr returns value as a []uint64, otherwise the provided default
func (m *Param) UintSliceOr(def []uint64) []uint64 {
	arr, ok := m.UintSlice()
	if ok {
		return arr
	}
	return def
}

// Int returns value as a int64, if possible
func (m *Param) Int() (int64, bool) {
	return toInt(m.value)
}

// IntOr returns value as a int64, otherwise the provided default
func (m *Param) IntOr(def int64) int64 {
	i, ok := m.Int()
	if ok {
		return i
	}
	return def
}

// IntSlice returns value as a []int64, if possible
func (m *Param) IntSlice() ([]int64, bool) {
	if m.value == nil {
		return nil, false
	}
	var res []int64
	res, ok := m.value.([]int64)

	if ok {
		return res, ok
	}

	r, ok := toSliceOf(m.value, intType, untypedInt)
	if !ok {
		return nil, false
	}
	res, ok = r.([]int64)
	if ok {
		return res, true
	}
	return nil, false
}

// IntSliceOr returns value as a []int64, otherwise the provided default
func (m *Param) IntSliceOr(def []int64) []int64 {
	arr, ok := m.IntSlice()

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

	return toFloat(m.value)
}

// FloatOr returns value as a float64, otherwise the provided default
func (m *Param) FloatOr(def float64) float64 {
	i, ok := m.Float()
	if ok {
		return i
	}
	return def
}

// FloatSlice returns value as a []float64, if possible
func (m *Param) FloatSlice() ([]float64, bool) {
	if m.value == nil {
		return nil, false
	}

	var res []float64
	res, ok := m.value.([]float64)

	if ok {
		return res, ok
	}

	r, ok := toSliceOf(m.value, floatType, untypedFloat)
	if !ok {
		return nil, false
	}
	res, ok = r.([]float64)
	if ok {
		return res, true
	}
	return nil, false
}

// FloatSliceOr returns value as a []float64, otherwise the provided default
func (m *Param) FloatSliceOr(def []float64) []float64 {
	arr, ok := m.FloatSlice()

	if ok {
		return arr
	}
	return def
}

// Bool returns value as a bool, if possible
func (m *Param) Bool() (bool, bool) {
	return toBool(m.value)
}

// BoolOr returns value as a bool, otherwise the provided default
func (m *Param) BoolOr(def bool) bool {
	i, ok := m.Bool()
	if ok {
		return i
	}
	return def
}

// BoolSlice returns value as a []bool, if possible
func (m *Param) BoolSlice() ([]bool, bool) {
	if m.value == nil {
		return nil, false
	}

	var res []bool
	res, ok := m.value.([]bool)

	if ok {
		return res, ok
	}

	r, ok := toSliceOf(m.value, boolType, untypedBool)
	if !ok {
		return nil, false
	}
	res, ok = r.([]bool)
	if ok {
		return res, true
	}
	return nil, false
}

// BoolSliceOr returns value as a []bool, otherwise the provided default
func (m *Param) BoolSliceOr(def []bool) []bool {
	arr, ok := m.BoolSlice()

	if ok {
		return arr
	}
	return def
}

func toString(in interface{}) (res string, ok bool) {
	if in == nil {
		return "", false
	}

	switch in.(type) {
	case string:
		res, ok = in.(string)
	case []byte:
		var b []byte
		b, ok = in.([]byte)
		res = string(b)
	case []rune:
		var r []rune
		r, ok = in.([]rune)
		res = string(r)
	}
	return
}

func toBool(in interface{}) (res bool, ok bool) {
	if in == nil {
		return false, false
	}

	switch in.(type) {
	case bool:
		res, ok = in.(bool)
	}
	return
}

func toUint(num interface{}) (uint64, bool) {

	if num == nil {
		return 0, false
	}

	var i uint64
	ok := false

	switch num.(type) {
	case int, int8, int16, int32, int64:
		a := reflect.ValueOf(num).Int() // a has type int64
		return uint64(a), true
	case uint, uint8, uint16, uint32, uint64:
		a := reflect.ValueOf(num).Uint() // a has type uint64
		return a, true
	case float64:
		f, ok := num.(float64)
		return uint64(f), ok
	case float32:
		f, ok := num.(float32)
		return uint64(f), ok
	}

	return i, ok

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

func toFloat(num interface{}) (float64, bool) {

	if num == nil {
		return 0, false
	}

	var i float64
	ok := false

	// TODO maybe remove reflection
	switch num.(type) {
	case int, int8, int16, int32, int64:
		a := reflect.ValueOf(num).Int() // a has type int64
		return float64(a), true
	case uint, uint8, uint16, uint32, uint64:
		a := reflect.ValueOf(num).Uint() // a has type uint64
		return float64(a), true
	case float64:
		f, ok := num.(float64)
		return float64(f), ok
	case float32:
		f, ok := num.(float32)
		return float64(f), ok
	}

	return i, ok

}

func toSliceOf(value interface{}, typ reflect.Type, converter converter) (interface{}, bool) {

	if reflect.TypeOf(value).Kind() != reflect.Slice {
		return nil, false
	}
	slice := reflect.ValueOf(value)
	resSlice := reflect.MakeSlice(reflect.SliceOf(typ), slice.Len(), slice.Len())

	for i := 0; i < slice.Len(); i++ {

		val, ok := converter(slice.Index(i).Interface())
		if !ok {
			return nil, false
		}
		resSlice.Index(i).Set(reflect.ValueOf(val))
	}
	return resSlice.Interface(), true
}
