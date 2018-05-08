package simple

import (
	"bitbucket.org/modfin/yarf"
	"errors"
	"fmt"
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

	ms := req.Param("sleep").IntOr(1000)

	select {
	case <-time.After(time.Duration(ms) * time.Millisecond):
	case <-req.Ctx.Done():
		print(" Canceled at server,", req.Ctx.Err())
		return req.Ctx.Err()
	}

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

func cat(req *yarf.Msg, resp *yarf.Msg) (err error) {

	print(" Got request cat")
	arr := req.Param("arr").StringArrOr([]string{"No", "Data"})

	res := ""
	for _, item := range arr {
		res += item
	}

	resp.SetParam("res", res)

	return nil
}

// StartServer starts a integration server using provided yarf transport
func StartServer(serverTransport yarf.Transporter, verbose bool) {
	doPrint = verbose

	print("Creating server")
	server := yarf.NewServer(serverTransport, "a", "integration")

	print("Adding err handler")
	server.Handle("err", err)

	print("Adding rpc err handler")
	server.Handle("rpc-err", rpcErr)

	print("Adding add handler")
	server.Handle("add", add)

	server.Handle("cat", cat)

	server.Handle("sleep", sleep)

	print("Adding sub handler")
	server.Handle("sub", func(req *yarf.Msg, resp *yarf.Msg) (err error) {
		print(" Got request Sub")
		var t Tuple
		err = req.Bind(&t)
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
