package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ugorji/go/codec"
	"github.com/vaughan0/go-zmq"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type HandshakeResponse struct {
	EventType int     `codec:0`
	Timestamp float64 `codec:1`
	Payload   struct {
		ServerVersion      string `codec:"rbkit_server_version"`
		ProtocolVersion    string `codec:"rbkit_protocol_version"`
		ProcessName        string `codec:"process_name"`
		Pwd                string `codec:"pwd"`
		Pid                int    `codec:"pid"`
		ObjectTraceEnabled int    `codec:"object_trace_enabled"`
	} `codec:2`
}

var h codec.MsgpackHandle

func main() {
	h.WriteExt = true
	h.RawToString = true
	//var msg map[int]interface{}
	msg := HandshakeResponse{}

	ctx, err := zmq.NewContext()
	checkError(err)
	defer ctx.Close()

	commandSock, err := ctx.Socket(zmq.Req)
	checkError(err)
	defer commandSock.Close()

	err = commandSock.Connect("tcp://127.0.0.1:5556")
	checkError(err)
	//var h codec.Handle = new(codec.MsgpackHandle)

	for {
		fmt.Println("Triggerring GC")
		err = commandSock.Send([][]byte{[]byte("handshake")})
		checkError(err)
		fmt.Println("Triggerred GC")

		fmt.Println("waiting for data")
		parts, err := commandSock.Recv()
		checkError(err)
		joinedBytes := bytes.Join(parts, nil)

		if len(joinedBytes) != 2 {
			// This is a msgpack message
			for _, part := range parts {
				var dec *codec.Decoder = codec.NewDecoderBytes(part, &h)
				err = dec.Decode(&msg)
				checkError(err)
				fmt.Println(msg)
			}

		} else {
			fmt.Println("received ", string(joinedBytes))
		}

		time.Sleep(time.Second * 1)
	}
}
