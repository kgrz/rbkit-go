package main

import (
	"fmt"
	"time"

	"github.com/vaughan0/go-zmq"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

var v1 string

func main() {
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
		err = commandSock.Send([][]byte{[]byte("trigger_gc")})
		checkError(err)
		fmt.Println("Triggerred GC")

		fmt.Println("waiting for data")
		parts, err := commandSock.Recv()
		checkError(err)
		fmt.Println(parts)

		//joinedBytes := bytes.Join(parts, nil)

		//var dec *codec.Decoder = codec.NewDecoderBytes(joinedBytes, h)
		//err = dec.Decode(&v1)
		//checkError(err)

		//fmt.Println(v1)

		time.Sleep(time.Second * 1)
	}
}
