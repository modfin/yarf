package simple

import (
	"bitbucket.org/modfin/yarf"
	"fmt"
)

func printPre(text string) func(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {

	return func(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {

		fmt.Println(text)

		err := next()

		return err
	}

}

func printPost(text string) func(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {

	return func(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {

		err := next()
		fmt.Println(text)

		return err
	}

}
