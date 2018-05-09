package tnats

import (
	"context"
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/nuid"
)

const ctrlHeaderLen = 30
const ctrlHeaderPrefix = "_Y_CTRL."
const ctrlCancel = "CANCEL"
const cmdUpgrade = "UPGRADE"

type txrx struct {
	transporter *NatsTransporter
	upgraded    bool
	init        bool

	function string
	message  *nats.Msg

	tx   string
	rx   string
	ctrl string
}

func (n *NatsTransporter) fromFunction(function string) txrx {

	// TODO add optional if context is provided
	ctrl := "_Y_CTRL." + nuid.Next()

	return txrx{
		transporter: n,
		upgraded:    false,
		function:    function,
		ctrl:        ctrl,
		init:        true,
	}
}

func (n *NatsTransporter) fromMessage(message *nats.Msg) txrx {

	// TODO add optional if context is provided
	ctrl := string(message.Data[:ctrlHeaderLen])
	message.Data = message.Data[ctrlHeaderLen:]

	return txrx{
		transporter: n,
		upgraded:    false,
		function:    message.Reply,
		message:     message,
		ctrl:        ctrl,
		init:        false,
	}
}

func (t *txrx) contextCanceler() (context.Context, func()) {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sub, err := t.transporter.client.SubscribeSync(t.ctrl)
		defer sub.Unsubscribe()
		defer cancel()

		if err != nil {
			fmt.Println("contextCanceler err 1", err)
		}

		for {
			msg, err := sub.NextMsgWithContext(ctx)
			// Context canceled
			if err != nil {
				//fmt.Println("HAS BIN CANCELD")
				return
			}

			//Canceling context
			if string(msg.Data) == ctrlCancel {
				//fmt.Println("IM CANCELING")
				return
			}
		}
	}()

	return ctx, cancel
}

func (t *txrx) send(ctx context.Context, data []byte) (err error) {
	var prefix string
	if t.init {
		prefix = t.ctrl
		t.init = false
	}

	if t.upgraded {
		return t.transporter.sendMultipart(t.tx, data)
	}

	if int(t.transporter.client.MaxPayload()) < len(data) {
		t.tx, t.rx, err = t.transporter.requestUpgrade(t.function, prefix)
		t.upgraded = true
		if err != nil {
			return err
		}
		return t.send(ctx, data)
	}

	// Init connection or just reply
	if t.message == nil {
		t.message, err = t.transporter.client.RequestWithContext(ctx, t.function, append([]byte(prefix), data...))
	} else {
		err = t.transporter.client.Publish(t.message.Reply, data)
	}

	return

}

func (t *txrx) receive(ctx context.Context) (data []byte, err error) {
	if t.upgraded {
		return t.transporter.receiveMultipart(t.rx)
	}

	if len(t.message.Data) > len(cmdUpgrade) && string(t.message.Data[:len(cmdUpgrade)]) == cmdUpgrade {
		t.tx, t.rx, err = t.transporter.acceptUpgrade(t.message)
		t.upgraded = true
		if err != nil {
			return nil, err
		}
		return t.receive(ctx)
	}

	return t.message.Data, nil
}
