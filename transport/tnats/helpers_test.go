package tnats

import (
	"testing"
)

func TestMin(t *testing.T) {


	if min(1,2) != 1 {
		t.Fail()
	}

	if min(1,-2) != -2 {
		t.Fail()
	}

}


func TestIntToBytes(t *testing.T) {

	i := 123

	b := intToBytes(i)

	j := bytesToInt(b)

	if i != j {
		t.Fail()
	}

}
