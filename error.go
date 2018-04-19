package yarf

import "strconv"

// RPCError struct for yarf rpc calls
type RPCError struct {
	Status int
	Msg    string
}

// NewRPCError create a RPCError struct for yarf rpc calls
func NewRPCError(status int, msg string) RPCError {
	return RPCError{
		Status: status,
		Msg:    msg,
	}
}

func (e RPCError) Error() string {
	return strconv.Itoa(e.Status) + ": " + e.Msg
}
