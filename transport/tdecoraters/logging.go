package tdecoraters

import (
	"bitbucket.org/modfin/yarf"
	"context"
	"time"
	"fmt"
)



type clientLoggingTransporter struct{
	next      yarf.Transporter
}

func (t clientLoggingTransporter) Call(ctx context.Context, function string, requestData []byte) (response []byte, err error){


	start := time.Now()
	response, err =  t.next.Call(ctx, function, requestData)

	fmt.Println("[", start.Format(time.RFC3339), "] Calling", function, "size", len(requestData), "time", time.Now().Sub(start).Truncate(time.Microsecond) )

	return response, err
}
func (t clientLoggingTransporter) Listen(function string, toExec func(ctx context.Context, requestData []byte) (responseData []byte)) error{
	return t.next.Listen(function, toExec)
}


func ClientLogging(next yarf.Transporter) yarf.Transporter{
	return clientLoggingTransporter{
		next: next,
	}
}