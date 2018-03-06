package main

import (
	"bitbucket.org/modfin/yarf/transport/thttp"
	"bitbucket.org/modfin/yarf"
	"github.com/golang-plus/errors"
	"log"
	"fmt"
	"time"
)

type Tuple struct {
	Val1 int
	Val2 int
}


func main(){


	fmt.Println("Creating server transport")
	serverTransport, _ := thttp.NewHttpTransporter(thttp.Options{})

	fmt.Println("Creating server")
	server := yarf.NewServer(serverTransport)

	fmt.Println("Adding add handler")
	server.Handle("add", func(req *yarf.Msg, resp *yarf.Msg) (err error){

		fmt.Println("Got request Add")
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
	});


	fmt.Println("Adding sub handler")
	server.Handle("sub", func(req *yarf.Msg, resp *yarf.Msg) (err error){

		fmt.Println("Got request Sub")
		var t Tuple
		err = req.Bind(&t)
		if err != nil {
			return errors.New("could not bind to model")
		}

		resp.SetParam("res", t.Val1 - t.Val2)

		return nil
	});


	time.Sleep(200 * time.Millisecond)

	fmt.Println("Creating client transport")
	clientTransport, _ := thttp.NewHttpTransporter( thttp.Options{ Discovery: &thttp.DiscoveryDnsA{Host:"localhost"}})

	fmt.Println("Creating client")
	client := yarf.NewClient(clientTransport)

	fmt.Println("Performing request, add")
	res, err := client.Request("add").
		SetParam("val1", 5).
		SetParam("val2", 7).
		Exec().
		Get()

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Result of 5 + 7 =", res.Param("res").IntOr(-1))


	fmt.Println("Performing request, sub")
	res, err = client.Request("sub").
		Content(Tuple{32, 11}).
		Exec().
		Get()

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Result of 32 + 11 =", res.Param("res").IntOr(-1))


}