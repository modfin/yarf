package simple

import (
	"bitbucket.org/modfin/yarf"
	"context"
	"fmt"
	"log"
	"time"
)

// TimeoutRequest preforms a request to server to sleep, with context
func TimeoutRequest(ctx context.Context, client yarf.Client, sleepMS int) (err error) {
	_, err = client.Request("a.integration.sleep").
		WithContext(ctx).
		SetParam("sleep", sleepMS).
		Exec().
		Get()
	return err
}

// ErrorRequest shall return a simple error
func ErrorRequest(client yarf.Client) (err error) {
	tuple := Tuple{}
	_, err = client.Request("a.integration.err").
		Exec().
		Bind(&tuple).
		Get()
	return err
}

// ErrorRequest2 shall return a simple error
func ErrorRequest2(client yarf.Client) (err error) {
	_, err = client.Request("a.integration.rpc-err").
		Exec().
		Get()
	return err
}

// ErrorChannelRequest shall return a simple error by using channels
func Error2ChannelRequest(client yarf.Client) (err error) {
	msgChan, errChan := client.Request("a.integration.rpc-err").
		UseChannels().
		Exec().
		Channels()

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
		msgChan<-msg
	}
	errFunc :=	func(e error) {
			errChan<-e
		}

	client.Request("a.integration.rpc-err").
		WithCallback(msgFunc, errFunc).
		Exec()

	select {
	case <-msgChan:
	case err = <-errChan:
	}

	return err
}


// CatRequest concatenates an array of strings in a string
func CatRequest(client yarf.Client, arr ...string) (*yarf.Msg, error) {
	return client.Request("a.integration.cat").
		SetParam("arr", arr).
		Exec().
		Get()
}

// CatChannelRequest concatenates an array of strings in a string by using channels
func CatChannelRequest(client yarf.Client, arr ...string) (*yarf.Msg, error) {
	msgChan, errChan := client.Request("a.integration.cat").
		SetParam("arr", arr).
		UseChannels().
		Exec().
		Channels()

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
		msgChan<-msg
	}
	errFunc :=	func(e error) {
		errChan<-e
	}

	client.Request("a.integration.cat").
		SetParam("arr", arr).
		WithCallback(msgFunc, errFunc).
		Exec()

	var err error
	var msg *yarf.Msg
	select {
	case msg = <-msgChan:
	case err = <-errChan:
	}

	return msg, err
}



// AddRequest adds two numbers
func AddRequest(client yarf.Client, i, j int) (*yarf.Msg, error) {
	return client.Request("a.integration.add").
		SetParam("val1", i).
		SetParam("val2", j).
		Exec().
		Get()
}

// SubRequest subtracts j from i
func SubRequest(client yarf.Client, i, j int) (*yarf.Msg, error) {
	return client.Request("a.integration.sub").
		Content(Tuple{i, j}).
		Exec().
		Get()
}

// LenRequest returns the length of an array
func LenRequest(client yarf.Client, len int) (*yarf.Msg, error) {
	arr := make([]byte, len)
	return client.Request("a.integration.len").
		BinaryContent(arr).
		Exec().
		Get()
}

// GenRequest generates an empty array of len
func GenRequest(client yarf.Client, len int) (*yarf.Msg, error) {
	return client.Request("a.integration.gen").
		SetParam("len", len).
		Exec().
		Get()
}

// CopyRequest makes a copy of an array of len
func CopyRequest(client yarf.Client, len int) (*yarf.Msg, error) {
	arr := make([]byte, len)
	return client.Request("a.integration.copy").
		BinaryContent(arr).
		Exec().
		Get()
}

// RunClient uses provided transport interface to run some tests
func RunClient(clientTransport yarf.Transporter) {

	var err error

	fmt.Println("Creating client")
	client := yarf.NewClient(clientTransport)

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

	time.Sleep(2 * time.Second)

}
