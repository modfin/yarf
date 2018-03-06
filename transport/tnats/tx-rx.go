package tnats

import (
	"github.com/nats-io/go-nats"
	"context"
)

type txrx struct {
	transporter *NatsTransporter
	upgraded    bool

	function    string
	message     *nats.Msg

	tx          string
	rx          string
}

func (n *NatsTransporter) fromFunction(function string) txrx {
	return txrx{
		transporter: n,
		upgraded: false,
		function: function,
	}
}

func (n *NatsTransporter) fromMessage(message *nats.Msg) txrx {
	return txrx{
		transporter: n,
		upgraded: false,
		function: message.Reply,
		message: message,
	}
}

func (t *txrx) send(ctx context.Context, data []byte) (err error) {

	if t.upgraded {
		return t.transporter.sendMultipart(t.tx, data)
	}

	if int(t.transporter.client.MaxPayload()) < len(data) {
		t.tx, t.rx, err = t.transporter.requestUpgrade(t.function);
		t.upgraded = true
		if err != nil {
			return err
		}
		return t.send(ctx, data)
	}

	// Init connection or just reply
	if t.message == nil {
		t.message, err = t.transporter.client.RequestWithContext(ctx, t.function, data)
	} else {
		err = t.transporter.client.Publish(t.message.Reply, data)
	}

	return

}
func (t *txrx) receive(ctx context.Context) (data []byte, err error) {
	if t.upgraded {
		return t.transporter.receiveMultipart(t.rx)
	}

	if (len(t.message.Data) > 8 && string(t.message.Data[:7]) == "UPGRADE") {
		t.tx, t.rx, err = t.transporter.acceptUpgrade(t.message);
		t.upgraded = true
		if err != nil {
			return nil, err
		}
		return t.receive(ctx)
	}

	return t.message.Data, nil
}



