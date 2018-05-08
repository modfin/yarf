package integration

import (
	"testing"
	"bitbucket.org/modfin/yarf"
	"bitbucket.org/modfin/yarf/example/simple"
	"strings"
)

func GetIntegrationTest(client yarf.Client) (func(t *testing.T)) {

	len := 2 * 1000000 // 2 mb

	return func(t *testing.T) {
		t.Run("GetTestErrors", GetTestErrors(client))
		t.Run("GetTestErrors", GetTestErrors2(client))
		t.Run("GetTestCat", GetTestCat(client, "a", "b", "c"))
		t.Run("GetTestAdd", GetTestAdd(client, 5, 7))
		t.Run("GetTestSub", GetTestSub(client, 33, 11))
		t.Run("GetTestLen", GetTestLen(client, len))
		t.Run("GetTestGen", GetTestGen(client, len))
		t.Run("GetTestCopy", GetTestCopy(client, len))

	}
}

func GetTestErrors(client yarf.Client) func(t *testing.T) {
	return func(t *testing.T) {
		err := simple.ErrorRequest(client)

		if err == nil {
			t.Log("Did not get an error")
			t.Fail()
			return
		}

		rerr, ok := err.(yarf.RPCError)

		if !ok {
			t.Fail()
			return
		}

		if rerr.Status != 510 {
			t.Log("expected 510 got,", rerr.Status)
			t.Fail()
		}

	}
}

func GetTestErrors2(client yarf.Client) func(t *testing.T) {
	return func(t *testing.T) {
		err := simple.ErrorRequest2(client)

		if err == nil {
			t.Log("Did not get an error")
			t.Fail()
			return
		}

		rerr, ok := err.(yarf.RPCError)

		if !ok {
			t.Fail()
			return
		}

		if rerr.Status != 600 {
			t.Log("expected 600 got,", rerr.Status)
			t.Fail()
		}

	}
}


func GetTestCat(client yarf.Client, arr ... string) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.CatRequest(client, arr...)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
		}

		if strings.Join(arr, "") !=  res.Param("res").StringOr("FAIL"){
			t.Log("Got response",  res.Param("res"), "expected", strings.Join(arr, "") )
			t.Fail()
		}


	}
}

func GetBenchmarkAdd(client yarf.Client, i, j int) func(t *testing.B) {
	return func(b *testing.B) {

		for n := 0; n < b.N; n++ {

			res, err := simple.AddRequest(client, i, j)

			if err != nil {
				b.Log("Got err", err)
				b.Fail()
				return
			}

			if int64(i+j) != res.Param("res").IntOr(-1) {
				b.Log("Got response", res.Param("res"), "expected", int64(i+j))
				b.Fail()
			}
		}

	}
}


func GetTestAdd(client yarf.Client, i, j int) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.AddRequest(client, i, j)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if int64(i+j) !=  res.Param("res").IntOr(-1) {
			t.Log("Got response",  res.Param("res"), "expected", int64(i+j))
			t.Fail()
		}
	}
}




func GetTestSub(client yarf.Client, i, j int) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.SubRequest(client, i, j)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if int64(i-j) !=  res.Param("res").IntOr(-1) {
			t.Log("Got response",  res.Param("res"), "expected", int64(i-j))
			t.Fail()
		}


	}
}


func GetTestLen(client yarf.Client, length int) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.LenRequest(client, length)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if int64(length) !=  res.Param("res").IntOr(-1) {
			t.Log("Got response",  res.Param("res"), "expected", length)
			t.Fail()
		}


	}
}


func GetTestGen(client yarf.Client, length int) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.GenRequest(client, length)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if length != len(res.Content) {
			t.Log("Got response len",  len(res.Content), "expected", length)
			t.Fail()
		}


	}
}

func GetTestCopy(client yarf.Client, length int) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.CopyRequest(client, length)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if length != len(res.Content) {
			t.Log("Got response len",  len(res.Content), "expected", length)
			t.Fail()
		}


	}
}