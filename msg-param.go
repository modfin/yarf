package yarf

type Param struct {
	key   string
	value interface{}
}

func (m *Param) Key() string {
	return m.key
}

func (m *Param) Value() interface{} {
	return m.value
}

func (m *Param) String() (string, bool) {
	if m.value == nil {
		return "", false
	}

	str, ok := m.value.(string)
	return str, ok
}

func (m *Param) StringOr(def string) string {
	str, ok := m.String()
	if ok {
		return str
	}
	return def
}

func (m *Param) StringArr() ([]string, bool) {
	if m.value == nil {
		return nil, false
	}
	arr, ok := m.value.([]interface{})
	var res []string
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

func (m *Param) StringArrOr(def []string) []string {
	arr, ok := m.StringArr()

	if ok {
		return arr
	}
	return def
}

func (m *Param) Int() (int64, bool) {

	if m.value == nil {
		return 0, false
	}

	var i int64
	var f float64
	ok := false

	switch m.value.(type) {
	case int64:
		i, ok = m.value.(int64)
	case float64:
		f, ok = m.value.(float64)
		i = int64(f)
	}

	return i, ok
}

func (m *Param) IntOr(def int64) int64 {
	i, ok := m.Int()
	if ok {
		return i
	}
	return def
}

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

func (m *Param) IntArrOr(def []int64) []int64 {
	arr, ok := m.IntArr()

	if ok {
		return arr
	}
	return def
}

func (m *Param) Float() (float64, bool) {
	if m.value == nil {
		return 0.0, false
	}

	i, ok := m.value.(float64)
	return i, ok
}

func (m *Param) FloatOr(def float64) float64 {
	i, ok := m.Float()
	if ok {
		return i
	}
	return def
}

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

func (m *Param) FloatArrOr(def []float64) []float64 {
	arr, ok := m.FloatArr()

	if ok {
		return arr
	}
	return def
}

func (m *Param) Bool() (bool, bool) {
	if m.value == nil {
		return false, false
	}

	i, ok := m.value.(bool)
	return i, ok
}
func (m *Param) BoolOr(def bool) bool {
	i, ok := m.Bool()
	if ok {
		return i
	}
	return def
}

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

func (m *Param) BoolArrOr(def []bool) []bool {
	arr, ok := m.BoolArr()

	if ok {
		return arr
	}
	return def
}
