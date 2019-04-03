package tnats

import (
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/go-nats"
	"strings"
	"time"
)

// NatsTransporter implements the yarf.Transport for using Nats
type NatsTransporter struct {
	namespace string
	servers   string
	opts      []nats.Option

	timeout time.Duration
	client  *nats.Conn
}

// NewNatsTransporter a constructor for the NatsTransporter
func NewNatsTransporter(servers string, timeout time.Duration, opts ...nats.Option) (*NatsTransporter, error) {
	t := NatsTransporter{
		namespace: "yarf.",
		servers:   servers,
		timeout:   timeout,
		opts:      opts,
	}

	nc, err := nats.Connect(servers, opts...)

	if err != nil {
		return nil, err
	}
	t.client = nc

	return &t, nil
}

// NewNatsTransporterFromConn a constructor for the NatsTransporter using an existing nats connection
func NewNatsTransporterFromConn(natsConnection *nats.Conn, timeout time.Duration) (*NatsTransporter, error) {
	t := NatsTransporter{
		namespace: "yarf.",
		timeout:   timeout,
		client:    natsConnection,
	}
	if !natsConnection.IsConnected() && !natsConnection.IsReconnecting() {
		return nil, errors.New("existing nats connection unusable")
	}
	return &t, nil
}

// Close the nats transporter and the nats client
func (n *NatsTransporter) Close() error {
	n.client.Close()
	return nil
}

// Call implements client side call of transporter
func (n *NatsTransporter) Call(ctx context.Context, function string, requestData []byte) (response []byte, err error) {

	if n.client.IsClosed() {
		return nil, errors.New("nats transporter has been closed")
	}

	// TODO if "Did not get messages in time nats: timeout" context does not seam to be canceled correctly after timeout ....
	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	function = n.namespace + function
	com := n.fromFunction(function)

	go func() {
		select {
		case <-ctx.Done():
			n.client.Publish(com.ctrl, []byte("CANCEL"))
		}
	}()

	err = com.send(ctx, requestData)
	if err != nil {
		return nil, err
	}

	response, err = com.receive(ctx)

	return
}

// Listen defines the function that will handle yarf requests
func (n *NatsTransporter) Listen(function string, toExec func(ctx context.Context, requestData []byte) (responseData []byte)) error {

	queueGroup := function
	parts := strings.Split(function, ".")
	if len(parts) > 1 {
		serverNamespaces := parts[:len(parts)-1]
		queueGroup = strings.Join(serverNamespaces, ".")
	}

	function = n.namespace + function
	queueGroup = n.namespace + queueGroup

	_, err0 := n.client.QueueSubscribe(function, queueGroup, func(m *nats.Msg) {
		go func() {
			com := n.fromMessage(m)

			ctx, cancel := com.contextCanceler()
			defer cancel()

			requestData, err := com.receive(ctx)
			if err != nil {
				fmt.Println("Could not receive ", err)
				return
			}

			responseData := toExec(ctx, requestData)

			err = com.send(ctx, responseData)
			if err != nil {
				fmt.Println("Could not send ", err)
				return
			}
		}()
	})

	return err0
}
