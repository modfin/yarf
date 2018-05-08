package tnats

import (
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/go-nats"
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

// Call implements client side call of transporter
func (n *NatsTransporter) Call(ctx context.Context, function string, requestData []byte) (response []byte, err error) {
	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	function = n.namespace + function
	com := n.fromFunction(function)

	err = com.send(ctx, requestData)
	if err != nil {
		return nil, err
	}

	response, err = com.receive(ctx)

	return
}

// Listen defines the function that will handle yarf requests
func (n *NatsTransporter) Listen(function string, toExec func(ctx context.Context, requestData []byte) (responseData []byte)) error {

	function = n.namespace + function
	_, err := n.client.Subscribe(function, func(m *nats.Msg) {

		com := n.fromMessage(m)

		requestData, err := com.receive(context.Background())
		if err != nil {
			fmt.Println("Could not recive ", err)
			return
		}

		responseData := toExec(context.Background(), requestData)

		err = com.send(context.Background(), responseData)
		if err != nil {
			fmt.Println("Could not send ", err)
			return
		}

	})

	return err
}
