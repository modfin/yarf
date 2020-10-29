package simple

import (
	"github.com/modfin/yarf"
)

func setValsMiddleware(i, j int) func(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {

	return func(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {

		request.SetParam("val1", i)
		request.SetParam("val2", j)

		err := next()

		return err
	}

}

func doubleResMiddleware(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {
	err := next()

	if err != nil {
		return err
	}
	i := response.Param("res").IntOr(0)
	response.SetParam("res", i*2)
	return nil
}
