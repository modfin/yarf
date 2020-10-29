package simple

import (
	"github.com/modfin/yarf"
	"github.com/modfin/yarf/middleware"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"time"
)

// TimeoutRequest preforms a request to server to sleep, with context
func TimeoutRequest(ctx context.Context, client yarf.Client, sleepMS int) (err error) {
	_, err = client.Request("a.integration.sleep").
		WithContext(ctx).
		WithParam("sleep", sleepMS).
		Get()
	return err
}

// ErrorRequest shall return a simple error
func ErrorRequest(client yarf.Client) (err error) {
	tuple := Tuple{}
	err = client.Request("a.integration.err").
		BindResponseContent(&tuple).
		Done()
	return err
}

// ErrorRequest2 shall return a simple error
func ErrorRequest2(client yarf.Client) (err error) {
	_, err = client.Request("a.integration.rpc-err").Get()
	return err
}

// Error2ChannelRequest shall return a simple error by using channels
func Error2ChannelRequest(client yarf.Client) (err error) {
	msgChan, errChan := client.Request("a.integration.rpc-err").Channels()

	select {
	case <-msgChan:
	case err = <-errChan:
	}

	return err
}

// Error2CallbackRequest shall return a simple error by using callbacks
func Error2CallbackRequest(client yarf.Client) (err error) {

	msgChan := make(chan *yarf.Msg)
	errChan := make(chan error)

	msgFunc := func(msg *yarf.Msg) {
		msgChan <- msg
	}
	errFunc := func(e error) {
		errChan <- e
	}

	client.Request("a.integration.rpc-err").Callbacks(msgFunc, errFunc)

	select {
	case <-msgChan:
	case err = <-errChan:
	}

	return err
}

// PanicRequest shall panic at the server side
func PanicRequest(client yarf.Client) (err error) {
	_, err = client.Request("a.integration.panic").
		Get()
	return err
}

// CatRequest concatenates an array of strings in a string
func CatRequest(client yarf.Client, arr ...string) (*yarf.Msg, error) {
	return client.Request("a.integration.cat").
		WithParam("arr", arr).
		Get()
}

// SHA256Request sha256 data...
func SHA256Request(client yarf.Client, data []byte) (string, error) {
	msg, err := client.Request("a.integration.sha256").
		WithBinaryContent(data).
		Get()

	if err != nil {
		return "", err
	}

	hash, ok := msg.Param("hash").String()

	if !ok {
		return "", errors.New("Could not get hash param as string")
	}

	return hash, nil
}

// CatChannelRequest concatenates an array of strings in a string by using channels
func CatChannelRequest(client yarf.Client, arr ...string) (*yarf.Msg, error) {
	msgChan, errChan := client.Request("a.integration.cat").
		WithParam("arr", arr).
		Channels()

	var err error
	var msg *yarf.Msg
	select {
	case msg = <-msgChan:
	case err = <-errChan:
	}

	return msg, err

}

// CatLateChannelRequest concatenates an array of strings in a string by using channels after the request has been resolved
func CatLateChannelRequest(client yarf.Client, arr ...string) (*yarf.Msg, error) {
	transit := client.Request("a.integration.cat").
		WithParam("arr", arr).
		Async()

	transit.Done()

	msgChan, errChan := transit.Channels()

	var err error
	var msg *yarf.Msg
	select {
	case msg = <-msgChan:
	case err = <-errChan:
	}

	return msg, err

}

// CatCallbackRequest concatenates an array of strings in a string by using callbacks
func CatCallbackRequest(client yarf.Client, arr ...string) (*yarf.Msg, error) {

	msgChan := make(chan *yarf.Msg)
	errChan := make(chan error)

	msgFunc := func(msg *yarf.Msg) {
		msgChan <- msg
	}
	errFunc := func(e error) {
		errChan <- e
	}

	client.Request("a.integration.cat").
		WithParam("arr", arr).
		Callbacks(msgFunc, errFunc)

	var err error
	var msg *yarf.Msg
	select {
	case msg = <-msgChan:
	case err = <-errChan:
	}

	return msg, err
}

// CatLateCallbackRequest concatenates an array of strings in a string by using callbacks after request is resolved
func CatLateCallbackRequest(client yarf.Client, arr ...string) (*yarf.Msg, error) {

	msgChan := make(chan *yarf.Msg)
	errChan := make(chan error)

	msgFunc := func(msg *yarf.Msg) {
		msgChan <- msg
	}
	errFunc := func(e error) {
		errChan <- e
	}

	transit := client.Request("a.integration.cat").
		WithParam("arr", arr).
		Async()

	transit.Done()

	transit.Callbacks(msgFunc, errFunc)

	var err error
	var msg *yarf.Msg
	select {
	case msg = <-msgChan:
	case err = <-errChan:
	}

	return msg, err
}

// SumRequest sums an integer array
func SumRequest(client yarf.Client, arr []int) (*yarf.Msg, error) {
	return client.Request("a.integration.sum").
		WithParam("arr", arr).
		Get()
}

// SumFloatRequest sums an float array
func SumFloatRequest(client yarf.Client, arr []float64) (*yarf.Msg, error) {
	return client.Request("a.integration.sumFloat").
		WithParam("arr", arr).
		Get()
}

// SumFloat32Request sums an float array
func SumFloat32Request(client yarf.Client, arr []float32) (*yarf.Msg, error) {
	return client.Request("a.integration.sumFloat32").
		WithParam("arr", arr).
		Get()
}

// XORRequest xors 2 arrays
func XORRequest(client yarf.Client, arr0 []bool, arr1 []bool) (*yarf.Msg, error) {
	return client.Request("a.integration.xor").
		WithParam("arr0", arr0).
		WithParam("arr1", arr1).
		Get()
}

// AddFloat32Request adds two float numbers
func AddFloat32Request(client yarf.Client, i, j float32) (*yarf.Msg, error) {
	return client.Request("a.integration.addFloat32").
		WithParam("val1", i).
		WithParam("val2", j).
		Get()
}

// AddFloatRequest adds two float numbers
func AddFloatRequest(client yarf.Client, i, j float64) (*yarf.Msg, error) {
	return client.Request("a.integration.addFloat").
		WithParam("val1", i).
		WithParam("val2", j).
		Get()
}

// AddRequest adds two numbers
func AddRequest(client yarf.Client, i, j int) (*yarf.Msg, error) {
	return client.Request("a.integration.add").
		WithParam("val1", i).
		WithParam("val2", j).
		Get()
}

// AddObserversRequest preforms a request that multiple observes
func AddObserversRequest(client yarf.Client, i, j, observers int) int {
	req := client.Request("a.integration.add").
		WithParam("val1", i).
		WithParam("val2", j)

	getAndEmmit := func(request *yarf.RPC, channel chan<- int) {
		msg, err := request.Get()
		if err != nil {
			channel <- 0
		}
		channel <- int(msg.Param("res").IntOr(0))

	}

	sumChan := make(chan int)

	for m := 0; m < observers; m++ {
		go getAndEmmit(req, sumChan)
	}

	sum := 0
	for m := 0; m < observers; m++ {
		select {
		case n := <-sumChan:
			sum += n
		}
	}

	return sum
}

// AddAndDoubleWithMiddlewareRequest adds two numbers and doubling result by using middleware
func AddAndDoubleWithMiddlewareRequest(client yarf.Client, i, j int) (*yarf.Msg, error) {
	return client.Request("a.integration.add").
		WithMiddleware(setValsMiddleware(i, j), doubleResMiddleware).
		Get()
}

// SubRequest subtracts j from i
func SubRequest(client yarf.Client, i, j int) (*yarf.Msg, error) {
	return client.Request("a.integration.sub").
		WithContent(Tuple{i, j}).
		Get()
}

// SwapAndMultiplyRequest swap places on values in tuple
func SwapAndMultiplyRequest(client yarf.Client, tuple Tuple, multiplier int) (res Tuple, err error) {

	err = client.Call(
		"a.integration.swapAndMultiply",
		tuple,
		&res,
		yarf.NewParam("multiplier", multiplier),
		yarf.NewParam("Something else", nil),
	)
	return
}

// SwapWithSerializer swap places on values in tuple
func SwapWithSerializer(client yarf.Client, tuple Tuple) (res Tuple, err error) {
	return res, client.Request("a.integration.swapWithSerializer").
		WithContentUsing(tuple, yarf.SerializerJson()).
		BindResponseContent(&res).
		Done()
}

// LenRequest returns the length of an array
func LenRequest(client yarf.Client, len int) (*yarf.Msg, error) {
	arr := make([]byte, len)
	return client.Request("a.integration.len").
		WithBinaryContent(arr).
		Get()
}

// GenRequest generates an empty array of len
func GenRequest(client yarf.Client, len int) (*yarf.Msg, error) {
	return client.Request("a.integration.gen").
		WithParam("len", len).
		Get()
}

// CopyRequest makes a copy of an array of len
func CopyRequest(client yarf.Client, len int) (*yarf.Msg, error) {
	arr := make([]byte, len)
	return client.Request("a.integration.copy").
		WithBinaryContent(arr).
		Get()
}



// GenRequest generates an empty array of len
func ConcRequest(client yarf.Client, sleep int) error {

	start := time.Now()
	m1, err1 := client.Request("a.integration.sleep").
		WithParam("sleep", sleep).
		Channels()

	m2, err2 := client.Request("a.integration.sleep").
		WithParam("sleep", sleep).
		Channels()

	st1 := time.Now()
	var err  error
	for i := 0; i < 2; i++ {
		select {
		case <-m1:
		case <-m2:
		case err = <-err1:
			break
		case err = <-err2:
			break
		}
	}

	if err != nil {
		return err
	}

	dur := time.Now().Sub(start)

	if dur > time.Duration(sleep*2)*time.Millisecond {
		err = fmt.Errorf("request did not occure concurrently, should have taken about %dms, took %s, st1 took %s", sleep, dur, st1.Sub(start))
	}

	return err
}



// RunClient uses provided transport interface to run some tests
func RunClient(clientTransport yarf.Transporter) {

	var err error

	fmt.Println("Creating client")
	client := yarf.NewClient(clientTransport)

	client.WithMiddleware(middleware.OpenTracing("Client> "))

	//client.WithMiddleware(PrintPre("Client 1"), PrintPre("Client 2"), PrintPre("Client 3"))

	fmt.Println("Performing timeout, sleep")
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	err = TimeoutRequest(ctx, client, 2000)
	cancel()
	fmt.Println(" Result of error", err)

	time.Sleep(1 * time.Second)

	fmt.Println("Performing request, err")
	err = ErrorRequest(client)
	fmt.Println(" Result of error", err)

	fmt.Println("Performing request, rpc-err")
	err = ErrorRequest2(client)
	fmt.Println(" Result of rpc error", err)

	fmt.Println("Performing request, panic")
	err = PanicRequest(client)
	fmt.Println(" Result of panic", err)

	var res *yarf.Msg
	fmt.Println("Performing request, cat")
	res, err = CatRequest(client, "a", "b", "c")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" Result of a + b + c =", res.Param("res").StringOr("Fail"))

	fmt.Println("Performing request, add")
	res, err = AddRequest(client, 5, 7)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" Result of 5 + 7 =", res.Param("res").IntOr(-1))

	fmt.Println("Performing request, sub")
	res, err = SubRequest(client, 32, 11)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(" Result of 32 + 11 =", res.Param("res").IntOr(-1))

	fmt.Println("Performing request, len")
	l := 2 * 1000000
	res, err = LenRequest(client, l)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(" Result ", l, "=", res.Param("res").IntOr(-1))

	fmt.Println("Performing request, gen")
	l = 2 * 1000000
	res, err = GenRequest(client, l)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(" Result  ", l, "=", len(res.Content))

	fmt.Println("Performing request, copy")
	l = 2 * 1000000
	res, err = CopyRequest(client, l)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(" Result ", l, "=", len(res.Content))

}
