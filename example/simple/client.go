package simple

import (
	"bitbucket.org/modfin/yarf"
	"fmt"
	"log"
)

func RunClinet(clientTransport yarf.Transporter){

	fmt.Println("Creating client")
	client := yarf.NewClient(clientTransport)




	fmt.Println("Performing request, cat")
	res, err := client.Request("a.test.cat").
		SetParam("arr", []string{"a","b", "c"}).
		Exec().
		Get()

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println(" Result of a + b + c =", res.Param("res").StringOr("Fail"))




	fmt.Println("Performing request, add")
	res, err = client.Request("a.test.add").
		SetParam("val1", 5).
		SetParam("val2", 7).
		Exec().
		Get()

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println(" Result of 5 + 7 =", res.Param("res").IntOr(-1))


	fmt.Println("Performing request, sub")
	res, err = client.Request("a.test.sub").
		Content(Tuple{32, 11}).
		Exec().
		Get()

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println(" Result of 32 + 11 =", res.Param("res").IntOr(-1))



	fmt.Println("Performing request, len")
	l := 2*1000000
	arr := make([]byte, l)
	res, err = client.Request("a.test.len").
		BinaryContent(arr).
		Exec().
		Get()

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println(" Result ", len(arr), "=", res.Param("res").IntOr(-1))



	fmt.Println("Performing request, gen")
	l = 2*1000000
	res, err = client.Request("a.test.gen").
		SetParam("len", l).
		Exec().
		Get()

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println(" Result  ", l, "=", len(res.Content))


	fmt.Println("Performing request, copy")
	l = 2*1000000
	arr = make([]byte, l)
	res, err = client.Request("a.test.copy").
		BinaryContent(arr).
		Exec().
		Get()

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println(" Result ", len(arr), "=", len(res.Content))

}