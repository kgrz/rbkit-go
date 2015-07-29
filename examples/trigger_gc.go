/* This program simply triggers a GC every 1 second and reads and parses the
* response. Useful for quick debugging of the network communication.
 */
package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/kgrz/rbkit-go/hw"
	"github.com/ugorji/go/codec"
	"github.com/vaughan0/go-zmq"
)

type HandshakeResponse struct {
	EventType int64
	Timestamp float64
	Payload   handshakePayload
}

type handshakePayload struct {
	ServerVersion      string
	ProtocolVersion    string
	ProcessName        string
	Pwd                string
	Pid                uint64
	ObjectTraceEnabled bool
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

var h codec.MsgpackHandle

func main() {
	h.WriteExt = true
	h.RawToString = true

	ctx, err := zmq.NewContext()
	checkError(err)
	defer ctx.Close()

	commandSock, err := ctx.Socket(zmq.Req)
	checkError(err)
	defer commandSock.Close()

	dataSock, err := ctx.Socket(zmq.Sub)
	checkError(err)
	defer dataSock.Close()

	err = commandSock.Connect("tcp://127.0.0.1:5556")
	checkError(err)

	dataSock.Subscribe([]byte(""))
	err = dataSock.Connect("tcp://127.0.0.1:5555")
	checkError(err)

	chans := commandSock.Channels()
	defer chans.Close()

	dataChans := dataSock.Channels()
	defer dataChans.Close()

	// send the command once. Do this inside a go routine
	sendCommand(chans)

	for {
		select {
		case parts := <-chans.In():
			go func() {
				fmt.Println("Received response")
				processIncomingMessage(parts)
				sendCommand(chans)
			}()
		case err := <-chans.Errors():
			checkError(err)
		case _ = <-dataChans.In():
			go func() {
				fmt.Println("Got info")
			}()
		case err := <-dataChans.Errors():
			checkError(err)
		}
		time.Sleep(time.Second * 1)
	}
}

func sendCommand(chans *zmq.Channels) {
	go func() {
		fmt.Println("Triggerring GC")
		chans.Out() <- [][]byte{[]byte("handshake")}
		fmt.Println("Triggerred GC")
	}()
}

func processIncomingMessage(parts [][]byte) {
	joinedBytes := bytes.Join(parts, nil)

	if len(joinedBytes) != 2 {
		// This is a msgpack message
		for _, part := range parts {
			unpacked := hw.Handshake{}

			_, err := unpacked.UnmarshalMsg(part)
			checkError(err)

			fmt.Println(unpacked)
		}
	} else {
		fmt.Println("received ", string(joinedBytes))
	}
}

func processIncomingEvent(parts [][]byte) {
	var msg map[int]interface{}

	for _, part := range parts {
		var dec *codec.Decoder = codec.NewDecoderBytes(part, &h)
		err := dec.Decode(&msg)
		checkError(err)
		fmt.Println(msg)
	}
}
