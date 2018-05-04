package simple

import (
	"bitbucket.org/modfin/yarf"
	"errors"
	"fmt"
)

type Tuple struct {
	Val1 int
	Val2 int
}



func err(req *yarf.Msg, resp *yarf.Msg) (err error) {

	fmt.Println(" Got request Error request")
	return errors.New("this endpoint returns an error")
}

func rpcErr(req *yarf.Msg, resp *yarf.Msg) (err error) {

	fmt.Println(" Got request RPC Error request")
	return yarf.NewRPCError(600, "custom error status")
}


func add(req *yarf.Msg, resp *yarf.Msg) (err error) {

	fmt.Println(" Got request Add")
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

	fmt.Println(" Got request cat")
	arr := req.Param("arr").StringArrOr([]string{"No", "Data"})

	res := ""
	for _, item := range arr {
		res += item
	}

	resp.SetParam("res", res)

	return nil
}

// StartServer starts a test server using provided yarf transport
func StartServer(serverTransport yarf.Transporter) {

	fmt.Println("Creating server")
	server := yarf.NewServer(serverTransport, "a", "test")

	fmt.Println("Adding err handler")
	server.Handle("err", err)

	fmt.Println("Adding rpc err handler")
	server.Handle("rpc-err", rpcErr)

	fmt.Println("Adding add handler")
	server.Handle("add", add)

	server.Handle("cat", cat)

	fmt.Println("Adding sub handler")
	server.Handle("sub", func(req *yarf.Msg, resp *yarf.Msg) (err error) {
		fmt.Println(" Got request Sub")
		var t Tuple
		err = req.Bind(&t)
		if err != nil {
			return errors.New("could not bind to model")
		}

		resp.SetParam("res", t.Val1-t.Val2)

		return nil
	})

	fmt.Println("Adding len content check")
	server.Handle("len", func(req *yarf.Msg, resp *yarf.Msg) (err error) {
		fmt.Println(" Got request len")

		resp.SetParam("res", len(req.Content))

		return nil
	})

	fmt.Println("Adding gen")
	server.Handle("gen", func(req *yarf.Msg, resp *yarf.Msg) (err error) {
		fmt.Println(" Got request gen")

		l, ok := req.Param("len").Int()
		if !ok {
			return errors.New("Could not read param len")
		}

		arr := make([]byte, l)
		resp.SetBinaryContent(arr)

		return nil
	})

	fmt.Println("Adding copy")
	server.Handle("copy", func(req *yarf.Msg, resp *yarf.Msg) (err error) {
		fmt.Println(" Got request copy")

		resp.SetBinaryContent(req.Content)

		return nil
	})

}
