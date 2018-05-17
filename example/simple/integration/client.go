package integration

import (
	"bitbucket.org/modfin/yarf"
	"bitbucket.org/modfin/yarf/example/simple"
	"context"
	"strings"
	"testing"
	"time"
)

// GetIntegrationTest generates a integration test for a specific client
func GetIntegrationTest(client yarf.Client) func(t *testing.T) {

	length := 2 * 1000000 // 2 mb

	return func(t *testing.T) {
		t.Run("GetTestContextTimeout", GetTestContextTimeout(client))
		t.Run("GetTestErrors", GetTestErrors(client))
		t.Run("GetTestErrors2", GetTestErrors2(client, simple.ErrorRequest2))
		t.Run("GetTestErrors2Channel", GetTestErrors2(client, simple.Error2ChannelRequest))
		t.Run("GetTestErrors2Callback", GetTestErrors2(client, simple.Error2CallbackRequest))
		t.Run("GetTestPanic", GetTestPanic(client, simple.PanicRequest))
		t.Run("GetTestCat", GetTestCat(client, simple.CatRequest, "a", "b", "c"))
		t.Run("GetTestCatChannel", GetTestCat(client, simple.CatChannelRequest, "a", "b", "c"))
		t.Run("GetTestCatCallback", GetTestCat(client, simple.CatCallbackRequest, "a", "b", "c"))
		t.Run("GetTestAdd", GetTestAdd(client, 5, 7))
		t.Run("GetTestAddAndDoubleWithMiddleware", GetTestAddAndDoubleWithMiddleware(client, 5, 7))
		t.Run("GetTestSub", GetTestSub(client, 33, 11))
		t.Run("GetTestLen", GetTestLen(client, length))
		t.Run("GetTestGen", GetTestGen(client, length))
		t.Run("GetTestCopy", GetTestCopy(client, length))

	}
}

// GetTestContextTimeout generates a Error test for a specific client
func GetTestContextTimeout(client yarf.Client) func(t *testing.T) {
	return func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
		defer cancel()

		err := simple.TimeoutRequest(ctx, client, 3000)

		if err == nil {
			t.Log("Did not get an error")
			t.Fail()
			return
		}

		if !strings.HasSuffix(err.Error(), "context deadline exceeded") {
			t.Log("Expected context deadline exceeded, got ,", err)
			t.Fail()
		}

	}
}

// GetTestErrors generates a Error test for a specific client
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

// GetTestErrors2 generates a Error test for a specific client
func GetTestErrors2(client yarf.Client, function func(client yarf.Client) (err error)) func(t *testing.T) {
	return func(t *testing.T) {
		err := function(client)

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

// GetTestPanic generates a Error and test server side recover middleware
func GetTestPanic(client yarf.Client, function func(client yarf.Client) (err error)) func(t *testing.T) {
	return func(t *testing.T) {
		err := function(client)

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

		if rerr.Status != yarf.StatusInternalPanic {
			t.Log("expected 600 got,", rerr.Status)
			t.Fail()
		}

	}
}

// GetTestCat generates a Array param test for a specific client
func GetTestCat(client yarf.Client, function func(client yarf.Client, arr ...string) (*yarf.Msg, error), arr ...string) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := function(client, arr...)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
		}

		if strings.Join(arr, "") != res.Param("res").StringOr("FAIL") {
			t.Log("Got response", res.Param("res"), "expected", strings.Join(arr, ""))
			t.Fail()
		}

	}
}

// GetBenchmarkAdd generates a integer param benchmark for a specific client
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

// GetTestAddAndDoubleWithMiddleware adds two numbers and doubling result by using middleware
func GetTestAddAndDoubleWithMiddleware(client yarf.Client, i, j int) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.AddAndDoubleWithMiddlewareRequest(client, i, j)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if 2*int64(i+j) != res.Param("res").IntOr(-1) {
			t.Log("Got response", res.Param("res"), "expected", 2*int64(i+j))
			t.Fail()
		}
	}
}

// GetTestAdd generates a integer param test for a specific client
func GetTestAdd(client yarf.Client, i, j int) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.AddRequest(client, i, j)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if int64(i+j) != res.Param("res").IntOr(-1) {
			t.Log("Got response", res.Param("res"), "expected", int64(i+j))
			t.Fail()
		}
	}
}

// GetTestSub generates a integer param test for a specific client
func GetTestSub(client yarf.Client, i, j int) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.SubRequest(client, i, j)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if int64(i-j) != res.Param("res").IntOr(-1) {
			t.Log("Got response", res.Param("res"), "expected", int64(i-j))
			t.Fail()
		}

	}
}

// GetTestLen generates a large request payload test for a specific client
func GetTestLen(client yarf.Client, length int) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.LenRequest(client, length)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if int64(length) != res.Param("res").IntOr(-1) {
			t.Log("Got response", res.Param("res"), "expected", length)
			t.Fail()
		}

	}
}

// GetTestGen generates a large response payload test for a specific client
func GetTestGen(client yarf.Client, length int) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.GenRequest(client, length)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if length != len(res.Content) {
			t.Log("Got response len", len(res.Content), "expected", length)
			t.Fail()
		}

	}
}

// GetTestCopy generates a large request/response payload test for a specific client
func GetTestCopy(client yarf.Client, length int) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.CopyRequest(client, length)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if length != len(res.Content) {
			t.Log("Got response len", len(res.Content), "expected", length)
			t.Fail()
		}

	}
}
