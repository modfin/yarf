package yarf

import "testing"

func TestToInt(t *testing.T) {

	if i, ok := toInt(float64(3.2)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toInt(float32(3.9)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toInt(int(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toInt(int8(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toInt(int16(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toInt(int32(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toInt(int64(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toInt(uint8(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toInt(uint16(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toInt(uint32(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toInt(uint64(3)); !ok || i != 3 {
		t.Fail()
	}

}

func TestToFloat(t *testing.T) {

	if i, ok := toFloat(float64(3.2)); !ok || i != 3.2 {
		t.Fail()
	}

	f := float32(3.2)
	if i, ok := toFloat(f); !ok || float32(i) != f {
		t.Error("Expected", f, "got", i)
	}

	if i, ok := toFloat(int(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toFloat(int8(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toFloat(int16(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toFloat(int32(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toFloat(int64(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toFloat(uint8(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toFloat(uint16(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toFloat(uint32(3)); !ok || i != 3 {
		t.Fail()
	}

	if i, ok := toFloat(uint64(3)); !ok || i != 3 {
		t.Fail()
	}

}

func TestNewParam(t *testing.T) {

	p := NewParam("nil", nil)

	if !p.IsNil() {
		t.Error("Didnt get nil 1")
	}

	if p.Value() != nil {
		t.Error("Didnt get nil 2")
	}

	if p.Key() != "nil" {
		t.Error("Key is not nil ")
	}

	p = NewParam("float", 32.12)

	if f, ok := p.Float(); f != 32.12 || !ok {
		t.Error("Did not get correkt float ", ok, f)
	}

	p = NewParam("int", 32)

	if f, ok := p.Int(); f != 32 || !ok {
		t.Error("Did not get correkt int ", ok, f)
	}

}
