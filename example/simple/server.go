package simple

import (
	"github.com/modfin/yarf"
	"github.com/modfin/yarf/middleware"
	"github.com/modfin/yarf/serializers/jsoniterator"
	"github.com/modfin/yarf/serializers/msgpack"
	"context"
	hashing "crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"time"
)

// Tuple is a simple struct for a tuple
type Tuple struct {
	Val1 int
	Val2 int
}

var doPrint bool

func print(args ...interface{}) {
	if doPrint {
		fmt.Println(args...)
	}
}

func sleep(req *yarf.Msg, resp *yarf.Msg) (err error) {

	ctx := req.Context()
	ms := req.Param("sleep").IntOr(1000)

	ctx, cc := context.WithCancel(ctx)
	defer cc()

	select {
	case <-time.After(time.Duration(ms) * time.Millisecond):
		print(" Waited the entire time, BAD", req.Context().Err())
	case <-ctx.Done():
		print(" Canceled at server,", ctx.Err())
		return ctx.Err()
	}

	//fmt.Println(ctx.Err())

	resp.SetParam("res", ms)

	return nil
}

func err(req *yarf.Msg, resp *yarf.Msg) (err error) {
	print(" Got request Error request")
	return errors.New("this endpoint returns an error")
}

func rpcErr(req *yarf.Msg, resp *yarf.Msg) (err error) {

	print(" Got request RPC Error request")
	return yarf.NewRPCError(600, "custom error status")
}

func add(req *yarf.Msg, resp *yarf.Msg) (err error) {

	print(" Got request Add")
	val1 := req.Param("val1")
	val2 := req.Param("val2")

	v1, ok := val1.Int()
	if !ok {
		return errors.New("could v1 not cast to int")
	}

	v2, ok := val2.Int()
	if !ok {
		return errors.New("could v2 not cast to int")
	}

	resp.SetParam("res", v1+v2)

	return nil
}

func sum(req *yarf.Msg, resp *yarf.Msg) (err error) {

	print(" Got request sum")
	arr := req.Param("arr")

	v1, ok := arr.IntSlice()
	if !ok {
		return errors.New("could arr not cast to []int64")
	}

	var acc int64
	for _, n := range v1 {
		acc += n
	}

	resp.SetParam("res", acc)

	return nil
}

func xor(req *yarf.Msg, resp *yarf.Msg) (err error) {

	print(" Got request sum")
	arr0 := req.Param("arr0")
	arr1 := req.Param("arr1")

	x, ok := arr0.BoolSlice()
	if !ok {
		return errors.New("could arr not cast to []bool")
	}

	y, ok := arr1.BoolSlice()
	if !ok {
		return errors.New("could arr not cast to []bool")
	}

	if len(x) != len(y) {
		return errors.New("arrays are not the same size")
	}

	res := make([]bool, len(x))
	for i := range x {
		res[i] = (x[i] || y[i]) && !(x[i] && y[i])
	}

	resp.SetParam("res", res)

	return nil
}

func sumFloat(req *yarf.Msg, resp *yarf.Msg) (err error) {

	print(" Got request sum")
	arr := req.Param("arr")

	v1, ok := arr.FloatSlice()
	if !ok {
		return errors.New("could arr not cast to []float64")
	}

	var acc float64
	for _, n := range v1 {
		acc += n
	}

	resp.SetParam("res", acc)

	return nil
}

func sumFloat32(req *yarf.Msg, resp *yarf.Msg) (err error) {

	print(" Got request sum")
	arr := req.Param("arr")

	v1, ok := arr.FloatSlice()
	if !ok {
		return errors.New("could arr not cast to []float64")
	}

	var acc float32
	for _, n := range v1 {
		acc += float32(n)
	}

	resp.SetParam("res", acc)

	return nil
}

func addFloat(req *yarf.Msg, resp *yarf.Msg) (err error) {

	print(" Got request Add")
	val1 := req.Param("val1")
	val2 := req.Param("val2")

	v1, ok := val1.Float()
	if !ok {
		return errors.New("could v1 not cast to int")
	}

	v2, ok := val2.Float()
	if !ok {
		return errors.New("could v2 not cast to int")
	}

	resp.SetParam("res", v1+v2)

	return nil
}

func addFloat32(req *yarf.Msg, resp *yarf.Msg) (err error) {

	print(" Got request Add")
	val1 := req.Param("val1")
	val2 := req.Param("val2")

	v1, ok := val1.Float()
	if !ok {
		return errors.New("could v1 not cast to int")
	}

	v2, ok := val2.Float()
	if !ok {
		return errors.New("could v2 not cast to int")
	}

	resp.SetParam("res", v1+v2)

	return nil
}

func swapAndMultiply(req *yarf.Msg, resp *yarf.Msg) (err error) {

	t := Tuple{}

	multiplier := int(req.Param("multiplier").IntOr(1))
	err = req.BindContent(&t)

	tmp := t.Val1
	t.Val1 = t.Val2 * multiplier
	t.Val2 = tmp * multiplier
	resp.SetContent(t)

	return
}

func sha256(req *yarf.Msg, resp *yarf.Msg) (err error) {

	hash := hashing.Sum256(req.Content)
	resp.SetParam("hash",
		base64.StdEncoding.EncodeToString(hash[:]))
	return
}

func swapWithSerializer(req *yarf.Msg, resp *yarf.Msg) (err error) {

	t := Tuple{}
	err = req.BindContent(&t)

	tmp := t.Val1
	t.Val1 = t.Val2
	t.Val2 = tmp
	resp.SetContentUsing(t, jsoniterator.Serializer())

	return
}

func cat(req *yarf.Msg, resp *yarf.Msg) (err error) {

	print(" Got request cat", reflect.TypeOf(req.Param("arr").Value()))

	arr := req.Param("arr").StringSliceOr([]string{"No", "Data"})

	res := ""
	for _, item := range arr {
		res += item
	}

	resp.SetParam("res", res)

	return nil
}

// StartServer starts a integration server using provided yarf transport
func StartServer(serverTransport yarf.Transporter, verbose bool) {
	StartServerWithSerializer(serverTransport, verbose, msgpack.Serializer())
}

// StartServerWithSerializer starts a integration server using provided yarf transport and a specific Serializer
func StartServerWithSerializer(serverTransport yarf.Transporter, verbose bool, serializer yarf.Serializer, midware ...yarf.Middleware) {
	doPrint = verbose

	print("Creating server")
	server := yarf.NewServer(serverTransport, "a", "integration")
	server.WithProtocolSerializer(serializer)
	server.WithSerializer(serializer)
	server.WithMiddleware(midware...)

	print("Adding handler by func name")
	server.HandleFunc(err)
	server.HandleFunc(add)
	server.HandleFunc(addFloat)
	server.HandleFunc(addFloat32)
	server.HandleFunc(xor)
	server.HandleFunc(sum)
	server.HandleFunc(sumFloat)
	server.HandleFunc(sumFloat32)
	server.HandleFunc(cat)
	server.HandleFunc(sleep)
	server.HandleFunc(swapAndMultiply)
	server.HandleFunc(swapWithSerializer)
	server.HandleFunc(sha256)

	print("Adding rpc err handler")
	server.Handle("rpc-err", rpcErr)

	print("Adding panic handler")
	server.Handle("panic", func(req *yarf.Msg, resp *yarf.Msg) (err error) {
		panic("im suppose to panic")
	}, middleware.Recover)

	print("Adding sub handler")
	server.Handle("sub", func(req *yarf.Msg, resp *yarf.Msg) (err error) {
		print(" Got request Sub")
		var t Tuple
		err = req.BindContent(&t)
		if err != nil {
			return errors.New("could not bind to model")
		}

		resp.SetParam("res", t.Val1-t.Val2)

		return nil
	})

	print("Adding len content check")
	server.Handle("len", func(req *yarf.Msg, resp *yarf.Msg) (err error) {
		print(" Got request len")

		resp.SetParam("res", len(req.Content))

		return nil
	})

	print("Adding gen")
	server.Handle("gen", func(req *yarf.Msg, resp *yarf.Msg) (err error) {
		print(" Got request gen")

		l, ok := req.Param("len").Int()
		if !ok {
			return errors.New("Could not read param len")
		}

		arr := make([]byte, l)
		resp.SetBinaryContent(arr)

		return nil
	})

	print("Adding copy")
	server.Handle("copy", func(req *yarf.Msg, resp *yarf.Msg) (err error) {
		print(" Got request copy")

		resp.SetBinaryContent(req.Content)

		return nil
	})

}
