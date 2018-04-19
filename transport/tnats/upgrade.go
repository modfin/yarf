package tnats

import (
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/satori/go.uuid"
	"strings"
)

const HeaderSize = 4 * 5

func (n *NatsTransporter) acceptUpgrade(m *nats.Msg) (tx string, rx string, err error) {

	cmd := strings.Split(string(m.Data), " ")
	if len(cmd) != 2 {
		fmt.Println("Fail do upgrade, All headers was not provided")
		return
	}

	//"UPGRADE " + u4.String()
	rx = cmd[1] + "-req"
	tx = cmd[1] + "-resp"

	err = n.client.Publish(m.Reply, []byte("OK"))
	if err != nil {
		fmt.Println("Could not send OK")
		return
	}

	return
}

func (n *NatsTransporter) requestUpgrade(function string) (tx string, rx string, err error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return
	}

	uv4s := n.namespace + u4.String()

	upgradeRequest := "UPGRADE " + uv4s

	tx = uv4s + "-req"
	rx = uv4s + "-resp"

	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()
	msg, err := n.client.RequestWithContext(ctx, function, []byte(upgradeRequest))
	if err != nil {
		return tx, rx, err
	}

	if string(msg.Data) != "OK" {
		return tx, rx, errors.New("Did not revcive ok to upgrade")
	}

	return
}

func (n *NatsTransporter) sendMultipart(channel string, data []byte) (err error) {

	var payloadSize = int(n.client.MaxPayload())
	var totalLen = len(data)
	contentLen := payloadSize - HeaderSize
	frames := totalLen/contentLen + 1

	waitChan := make(chan bool)
	errorChan := make(chan error)

	for frame := 0; frame < frames; frame++ {

		start := frame * contentLen
		end := min(start+contentLen, totalLen)

		headers := make([]byte, 0, HeaderSize)
		headers = append(headers, intToBytes(totalLen)...)
		headers = append(headers, intToBytes(start)...)
		headers = append(headers, intToBytes(end)...)
		headers = append(headers, intToBytes(frame)...)
		headers = append(headers, intToBytes(frames)...)

		go func() {
			err = n.client.Publish(
				channel,
				append(headers, data[start:end]...))
			if err != nil {
				errorChan <- err
			} else {
				waitChan <- true
			}
		}()

		if err != nil {
			fmt.Println(1, err)
			return err
		}
	}

	for frame := 0; frame < frames; frame++ {
		select {
		case <-waitChan:
		case <-errorChan:
			if err != nil {
				fmt.Println(2, err)
				return err
			}
		}
	}

	return nil
}

func (n *NatsTransporter) receiveMultipart(channel string) (data []byte, err error) {

	sub, err := n.client.SubscribeSync(channel)
	defer func() {
		err = sub.Unsubscribe()
		if err != nil {
			fmt.Println("Could not unsubscribe")
		}

	}()
	if err != nil {
		fmt.Println("Could not send Subscribe")
		return
	}

	waitChan := make(chan bool, 2)
	errorsChan := make(chan error)

	frames := 1
	for i := 0; i < frames; i++ {

		resv := func() {
			msg, err := sub.NextMsg(n.timeout)
			if err != nil {
				fmt.Println("Did not get messages in time", err)
				errorsChan <- err
				return
			}

			partial := msg.Data

			// Extracting headers
			totalLen := bytesToInt(partial[0*4 : 0*4+4]) // append(headers, intToBytes(totalLen)...)
			start := bytesToInt(partial[1*4 : 1*4+4])    //append(headers, intToBytes(start)...)
			end := bytesToInt(partial[2*4 : 2*4+4])      //append(headers, intToBytes(end)...)
			//frame := bytesToInt(partial[3 * 4 : 3 * 4 + 4])//append(headers, intToBytes(frame)...)
			framesHeader := bytesToInt(partial[4*4 : 4*4+4]) //append(headers, intToBytes(frames)...)

			if data == nil {
				data = make([]byte, totalLen)
				frames = framesHeader
			}

			partial = partial[HeaderSize:]

			end = min(end, totalLen)

			copy(data[start:end], partial)

			waitChan <- true
		}

		// wait for first frame before trying to revice all
		if data == nil {
			resv()
		} else {
			go resv()
		}

	}

	for i := 0; i < frames; i++ {
		select {
		case <-waitChan:
		case err = <-errorsChan:
			return nil, err

		}
	}

	return data, err
}
