package yarf

import "strconv"

type RPCError struct {
	Status int
	Msg    string
}

func NewRPCError(status int, msg string) RPCError {
	return RPCError{
		Status: status,
		Msg:    msg,
	}
}

func (e RPCError) Error() string {
	return strconv.Itoa(e.Status) + ": " + e.Msg
}
