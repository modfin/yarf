package middleware

import (
	"bitbucket.org/modfin/yarf"
	"fmt"
)

// Recover recovers from panic and converts it to an error.
func Recover(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = yarf.NewRPCError(yarf.StatusInternalPanic, fmt.Sprintf("panic, %s", r))
		}
	}()

	err = next()

	return
}
