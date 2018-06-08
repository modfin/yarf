package yarf

import (
	"reflect"
	"testing"
)

var toIntTable = []struct {
	in       interface{}
	expected int64
}{
	{float64(3.2), 3},
	{float64(3.9), 3},
	{float32(3.2), 3},
	{byte(3), 3},
	{int(3), 3},
	{int8(3), 3},
	{int16(3), 3},
	{int32(3), 3},
	{int64(3), 3},

	{uint(3), 3},
	{uint8(3), 3},
	{uint16(3), 3},
	{uint32(3), 3},
	{uint64(3), 3},
}

func TestToInt(t *testing.T) {

	for _, f := range toIntTable {
		if i, ok := toInt(f.in); !ok || i != f.expected {
			t.Fail()
		}
	}

}

var toFloatTable = []struct {
	in       interface{}
	expected float64
}{
	{float64(3.2), 3.2},
	{float32(3.2), float64(float32(3.2))},
	{byte(3), 3},
	{int(3), 3},
	{int8(3), 3},
	{int16(3), 3},
	{int32(3), 3},
	{int64(3), 3},

	{uint(3), 3},
	{uint8(3), 3},
	{uint16(3), 3},
	{uint32(3), 3},
	{uint64(3), 3},
}

func TestToFloat(t *testing.T) {

	for _, f := range toFloatTable {
		if i, ok := toFloat(f.in); !ok || i != f.expected {
			t.Fail()
		}
	}

}

var newParamSimpleTable = []struct {
	key   string
	val   interface{}
	extra func(param Param) (ok bool, got interface{}, expected interface{})
}{
	{"float", 3.2, func(param Param) (ok bool, got interface{}, expected interface{}) {
		expected = 3.2
		got, ok = param.Float()
		return
	}},
	{"int", 3, func(param Param) (ok bool, got interface{}, expected interface{}) {
		expected = int64(3)
		got, ok = param.Int()
		return
	}},
	{"bool", true, func(param Param) (ok bool, got interface{}, expected interface{}) {
		expected = true
		got, ok = param.Bool()
		return
	}},
	{"bool false", false, func(param Param) (ok bool, got interface{}, expected interface{}) {
		expected = false
		got, ok = param.Bool()
		return
	}},
	{"string", "a string", func(param Param) (ok bool, got interface{}, expected interface{}) {
		expected = "a string"
		got, ok = param.String()
		return
	}},
	{"nil", nil, nil},
}

func TestNewSimpleParam(t *testing.T) {

	for _, f := range newParamSimpleTable {

		p := NewParam(f.key, f.val)

		if p.IsNil() != (f.val == nil) {
			t.Error("got", p.IsNil(), "expected", f.val == nil)
		}

		if p.IsSlice() {
			t.Error("value", f.val, "is not a slice")
		}

		if p.Value() != f.val {
			t.Error("got", p.Value(), "expected", f.val)

		}

		if p.Key() != f.key {
			t.Error("got", p.Key(), "expected", f.key)
		}

		if f.extra != nil {

			if ok, got, expected := f.extra(p); !ok || got != expected {
				t.Error("got", got, "expected", expected)
			}

		}
	}
}

var newParamSliceTable = []struct {
	key   string
	val   interface{}
	extra func(param Param) (ok bool, got interface{}, expected interface{})
}{
	{"float", []float64{3.2, 93.2}, func(param Param) (ok bool, got interface{}, expected interface{}) {
		expected = []float64{3.2, 93.2}
		got, ok = param.FloatSlice()
		return
	}},
	{"float32", []float32{3.2, 93.2}, func(param Param) (ok bool, got interface{}, expected interface{}) {
		expected = []float64{float64(float32(3.2)), float64(float32(93.2))}
		got, ok = param.FloatSlice()
		return
	}},

	{"int", []int64{3, 93}, func(param Param) (ok bool, got interface{}, expected interface{}) {
		expected = []int64{3, 93}
		got, ok = param.IntSlice()
		return
	}},
	{"uintcast", []uint{3, 93}, func(param Param) (ok bool, got interface{}, expected interface{}) {
		expected = []int64{3, 93}
		got, ok = param.IntSlice()
		return
	}},
	{"floatAsInt", []float64{3.2, 93.2}, func(param Param) (ok bool, got interface{}, expected interface{}) {
		expected = []int64{3, 93}
		got, ok = param.IntSlice()
		return
	}},
	{"intAsFloat", []int{3, 93}, func(param Param) (ok bool, got interface{}, expected interface{}) {
		expected = []float64{3, 93}
		got, ok = param.FloatSlice()
		return
	}},
	{"bool", []bool{true, false}, func(param Param) (ok bool, got interface{}, expected interface{}) {
		expected = []bool{true, false}
		got, ok = param.BoolSlice()
		return
	}},
}

func TestNewSliceParam(t *testing.T) {

	for _, f := range newParamSliceTable {

		p := NewParam(f.key, f.val)

		if p.Key() != f.key {
			t.Error("got", p.Key(), "expected", f.key)
		}

		if !p.IsSlice() {
			t.Error("value", f.val, "is a slice")
		}

		if p.IsNil() {
			t.Error("value", f.val, "is not nil")
		}

		if !reflect.DeepEqual(p.Value(), f.val) {
			t.Error("got", p.Value(), "expected", f.val)

		}

		if f.extra != nil {
			if ok, got, expected := f.extra(p); !ok || !reflect.DeepEqual(got, expected) {
				t.Error("got", got, "expected", expected)
			}
		}
	}
}
