package simple

import (
	"bitbucket.org/modfin/yarf"
	"fmt"
	"log"
)

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

// CatRequest concatenates an array of strings in a string
func CatRequest(client yarf.Client, arr ...string) (*yarf.Msg, error) {
	return client.Request("a.integration.cat").
		SetParam("arr", arr).
		Exec().
		Get()
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
	var res *yarf.Msg
	fmt.Println("Creating client")
	client := yarf.NewClient(clientTransport)

	fmt.Println("Performing request, err")
	err = ErrorRequest(client)
	fmt.Println(" Result of error", err)

	fmt.Println("Performing request, rpc-err")
	err = ErrorRequest2(client)
	fmt.Println(" Result of rpc error", err)

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
