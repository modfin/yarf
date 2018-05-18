package thttp

import (
	"testing"
	"time"
)

func TestStringOr(t *testing.T) {


	if stringOr("", "want") != "want"{
		t.Fail()
	}

	if stringOr("want", "not this") != "want"{
		t.Fail()
	}

}


func TestDurationOr(t *testing.T) {

	var d time.Duration

	if durationOr(d, time.Second) != time.Second{
		t.Fail()
	}

	if durationOr(time.Minute, time.Second) != time.Minute{
		t.Fail()
	}

}
