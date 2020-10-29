package integration

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/modfin/yarf"
	"github.com/modfin/yarf/example/simple"
	"strings"
	"testing"
	"time"
)

// GetIntegrationTest generates a integration test for a specific client
func GetIntegrationTest(client yarf.Client) func(t *testing.T) {

	return func(t *testing.T) {
		t.Run("ContextTimeout", GetTestContextTimeout(client))
		t.Run("Errors", GetTestErrors(client))
		t.Run("Errors2", GetTestErrors2(client, simple.ErrorRequest2))
		t.Run("Errors2Channel", GetTestErrors2(client, simple.Error2ChannelRequest))
		t.Run("Errors2Callback", GetTestErrors2(client, simple.Error2CallbackRequest))
		t.Run("Panic", GetTestPanic(client, simple.PanicRequest))
		t.Run("Cat", GetTestCat(client, simple.CatRequest, "a", "b", "c"))
		t.Run("CatChannel", GetTestCat(client, simple.CatChannelRequest, "a", "b", "c"))
		t.Run("CatLateChannel", GetTestCat(client, simple.CatLateChannelRequest, "a", "b", "c"))
		t.Run("CatCallback", GetTestCat(client, simple.CatCallbackRequest, "a", "b", "c"))
		t.Run("CatLateCallback", GetTestCat(client, simple.CatLateCallbackRequest, "a", "b", "c"))
		t.Run("XOR", GetTestXOR(client, []bool{true, true, true}, []bool{true, false, true}, []bool{false, true, false}))
		t.Run("Sum", GetTestSum(client, []int{3, 5, 7}))
		t.Run("SumFloat", GetTestSumFloat(client, []float64{3.2, 5.3, 7.4}))
		t.Run("SumFloat32", GetTestSumFloat32(client, []float32{3.1, 5.0, 7.1}))
		t.Run("Add", GetTestAdd(client, 5, 7))
		t.Run("AddFloat", GetTestAddFloat(client, 5.2, 7.3))
		t.Run("AddFloat32", GetTestAddFloat32(client, 5.2, 7.3))
		t.Run("AddAndDoubleWithMiddleware", GetTestAddAndDoubleWithMiddleware(client, 5, 7))
		t.Run("ObservedAdd", GetTestObservedAdd(client, 5, 7, 3))
		t.Run("Sub", GetTestSub(client, 33, 11))
		t.Run("Swap", GetTestSwap(client, simple.Tuple{Val1: 1, Val2: 2}, 3))
		t.Run("Conc", GetTestConc(client, 25))
		t.Run("SwapWithSerializer", GetTestSwapWithSerlizer(client, simple.Tuple{Val1: 1, Val2: 2}))
	}
}

// GetExtraIntegrationTest generates a integration test for a specific client with large payloads
func GetExtraIntegrationTest(client yarf.Client) func(t *testing.T) {
	length := 4 * 1000000 // ~4 mb
	return func(t *testing.T) {
		t.Run("GetTestLen", GetTestLen(client, length))
		t.Run("GetTestGen", GetTestGen(client, length))
		t.Run("GetTestCopy", GetTestCopy(client, length))
		t.Run("GetTestSHA256", GetTestSHA256(client, length))
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

// GetTestObservedAdd adds two numbers and multiplys it by the num of observers
func GetTestObservedAdd(client yarf.Client, i, j int, observers int) func(t *testing.T) {
	return func(t *testing.T) {
		res := simple.AddObserversRequest(client, i, j, observers)

		if observers*(i+j) != res {
			t.Log("Got response", res, "expected", observers*(i+j))
			t.Fail()
		}
	}
}

// GetTestAddFloat generates a float param test for a specific client
func GetTestAddFloat(client yarf.Client, i, j float64) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.AddFloatRequest(client, i, j)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if i+j != res.Param("res").FloatOr(-1.0) {
			t.Log("Got response", res.Param("res"), "expected", i+j)
			t.Fail()
		}
	}
}

// GetTestAddFloat32 generates a float param test for a specific client
func GetTestAddFloat32(client yarf.Client, i, j float32) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.AddFloat32Request(client, i, j)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if float64(i+j) != res.Param("res").FloatOr(-1.0) {
			t.Log("Got response", res.Param("res"), "expected", i+j)
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

// GetTestXOR sums a integer array
func GetTestXOR(client yarf.Client, arr0 []bool, arr1 []bool, expected []bool) func(t *testing.T) {
	return func(t *testing.T) {
		msg, err := simple.XORRequest(client, arr0, arr1)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		res, ok := msg.Param("res").BoolSlice()

		if !ok {
			t.Log("did not get an array")
			t.Fail()
			return
		}

		for i := range expected {
			if res[i] != expected[i] {
				t.Log("Got response", res, "expected", expected)
				t.Fail()
			}
		}

	}
}

// GetTestSum sums a integer array
func GetTestSum(client yarf.Client, arr []int) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.SumRequest(client, arr)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		var acc int
		for _, i := range arr {
			acc += i
		}

		if int64(acc) != res.Param("res").IntOr(-1) {
			t.Log("Got response", res.Param("res"), "expected", acc)
			t.Fail()
		}
	}
}

// GetTestSumFloat sums a float array
func GetTestSumFloat(client yarf.Client, arr []float64) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.SumFloatRequest(client, arr)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		var acc float64
		for _, i := range arr {
			acc += i
		}

		if acc != res.Param("res").FloatOr(-1.0) {
			t.Log("Got response", res.Param("res"), "expected", acc)
			t.Fail()
		}
	}
}

// GetTestSumFloat32 sums a float array
func GetTestSumFloat32(client yarf.Client, arr []float32) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.SumFloat32Request(client, arr)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		var acc float32
		for _, i := range arr {
			acc += i
		}

		got := float32(res.Param("res").FloatOr(-1.0))
		if acc != got {
			t.Log("Got response", got, "expected", acc)
			t.Fail()
		}
	}
}

// GetTestSwapWithSerlizer generates test for swaping values in a tuple using specific Json Sterilizer
func GetTestSwapWithSerlizer(client yarf.Client, tuple simple.Tuple) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.SwapWithSerializer(client, tuple)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if tuple.Val1 != res.Val2 && tuple.Val2 != res.Val1 {
			t.Log("Got response", res, "expected", simple.Tuple{Val1: tuple.Val2, Val2: tuple.Val1})
			t.Fail()
		}
	}
}

// GetTestSwap generates test for swaping values in a tuple
func GetTestSwap(client yarf.Client, tuple simple.Tuple, multiplier int) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := simple.SwapAndMultiplyRequest(client, tuple, multiplier)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}

		if tuple.Val1*multiplier != res.Val2 && tuple.Val2*multiplier != res.Val1 {
			t.Log("Got response", res, "expected", simple.Tuple{Val1: tuple.Val2, Val2: tuple.Val1})
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

// GetTestSHA256 generates 10 sha256 request up to len with random data
func GetTestSHA256(client yarf.Client, length int) func(t *testing.T) {
	return func(t *testing.T) {

		for i := 0; i < 10; i++ {

			data := make([]byte, length/(10-i))
			_, err := rand.Read(data)
			if err != nil {
				t.Log("Got err making random arr", err)
				t.Fail()
			}

			hash := sha256.Sum256(data)
			expected := base64.StdEncoding.EncodeToString(hash[:])

			got, err := simple.SHA256Request(client, data)

			if err != nil {
				t.Log("Got err from request", got)
				t.Fail()
			}

			if expected != got {
				t.Log("Expected", expected, "But got", got)
				t.Fail()
			}
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


// GetTestCopy generates a large request/response payload test for a specific client
func GetTestConc(client yarf.Client, sleep int) func(t *testing.T) {
	return func(t *testing.T) {
		err := simple.ConcRequest(client, sleep)

		if err != nil {
			t.Log("Got err", err)
			t.Fail()
			return
		}
	}
}