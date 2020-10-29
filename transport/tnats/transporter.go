package tnats

import (
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/go-nats"
	"strings"
	"sync"
	"time"
)

// NatsTransporter implements the yarf.Transport for using Nats
type NatsTransporter struct {
	namespace string
	servers   string
	opts      []nats.Option

	timeout time.Duration
	client  *nats.Conn

	mu     sync.Mutex
	count  int64
	subs   []*nats.Subscription
	closed chan struct{}
}

// NewNatsTransporter a constructor for the NatsTransporter
func NewNatsTransporter(servers string, timeout time.Duration, opts ...nats.Option) (*NatsTransporter, error) {
	t := NatsTransporter{
		namespace: "yarf.",
		servers:   servers,
		timeout:   timeout,
		opts:      opts,
		closed:    make(chan struct{}),
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
		closed:    make(chan struct{}),
	}
	if !natsConnection.IsConnected() && !natsConnection.IsReconnecting() {
		return nil, errors.New("existing nats connection unusable")
	}
	return &t, nil
}

func (n *NatsTransporter) IsClose() bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	select {
	case <-n.closed:
		return true
	default:
		return false
	}
}

// Close the nats transporter and the nats client
func (n *NatsTransporter) Close() error {
	n.mu.Lock()
	select {
	case <-n.closed:
	default:
		close(n.closed)
	}
	n.mu.Unlock()
	n.client.Close()
	return nil
}
func (n *NatsTransporter) incCount() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.count += 1
}
func (n *NatsTransporter) decCount() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.count -= 1
}

func (n *NatsTransporter) CloseGraceful(timeout time.Duration) error {
	timeoutChan := time.After(timeout)
	n.mu.Lock()
	select {
	case <-n.closed:
	default:
		close(n.closed)
	}
	var err error
	for _, s := range n.subs {
		err0 := s.Drain()
		if err0 != nil {
			err = err0
		}
	}
	n.mu.Unlock()
	// Waiting for things to finish
	for {
		n.mu.Lock()
		if n.count == 0 {
			n.mu.Unlock()
			break
		}
		n.mu.Unlock()
		select {
		case <-time.After(time.Millisecond * 50):
		case <-timeoutChan:
			break
		}
	}
	n.client.Close()
	return err
}

// Call implements client side call of transporter
func (n *NatsTransporter) Call(ctx context.Context, function string, requestData []byte) (response []byte, err error) {
	if n.IsClose() {
		return nil, errors.New("transport layer is has been closed")
	}
	n.incCount()
	defer n.decCount()

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

	if n.IsClose() {
		return errors.New("transport layer is has been closed")
	}

	queueGroup := function
	parts := strings.Split(function, ".")
	if len(parts) > 1 {
		serverNamespaces := parts[:len(parts)-1]
		queueGroup = strings.Join(serverNamespaces, ".")
	}

	function = n.namespace + function
	queueGroup = n.namespace + queueGroup

	sub, err0 := n.client.QueueSubscribe(function, queueGroup, func(m *nats.Msg) {
		go func() {
			n.incCount()
			defer n.decCount()
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

	n.mu.Lock()
	n.subs = append(n.subs, sub)
	n.mu.Unlock()

	return err0
}
