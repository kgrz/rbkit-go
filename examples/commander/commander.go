// Prompt user for commands, print them out. This has to be done in a
// separate thread
package main

import (
	"bytes"
	"fmt"

	"github.com/kgrz/rbkit-go/unpack"
	"github.com/ugorji/go/codec"
	"github.com/vaughan0/go-zmq"
)

var optionDict = map[int]string{
	1: "start_memory_profile",
	2: "stop_memory_profile",
	3: "objectspace_snapshot",
	4: "trigger_gc",
	5: "handshake",
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func askForOption() <-chan string {
	option := make(chan string)
	go func() {
		var input int
		fmt.Println("Enter a selection for the command: ")
		fmt.Println("1. Start Memory Profile")
		fmt.Println("2. Stop Memory Profile")
		fmt.Println("3. Objectspace Snapshot")
		fmt.Println("4. Trigger GC")
		fmt.Println("5. Handshake")
		fmt.Println("Hit Ctrl+C to stop")
		fmt.Scanln(&input)

		if input > 0 && input < 6 {
			option <- optionDict[input]
		} else {
			go askForOption()
			fmt.Println("\nInvalid option\n\n")
		}
	}()
	return option
}

func main() {
	ctx, err := zmq.NewContext()
	checkError(err)
	defer ctx.Close()

	commandSock, err := ctx.Socket(zmq.Req)
	checkError(err)
	defer commandSock.Close()

	err = commandSock.Connect("tcp://127.0.0.1:5556")
	checkError(err)
	chans := commandSock.Channels()

	// Repeatedly ask for an option
	for {
		optionChan := askForOption()
		chans.Out() <- [][]byte{[]byte(<-optionChan)}

		select {
		case parts := <-chans.In():
			decodeMsgpackAndPrint(parts)
		case err := <-chans.Errors():
			checkError(err)
		}
	}
}

func decodeMsgpackAndPrint(parts [][]byte) {
	var h codec.MsgpackHandle
	h.WriteExt = true
	h.RawToString = true
	var msg map[int]interface{}
	joinedBytes := bytes.Join(parts, nil)

	if len(joinedBytes) != 2 {
		// This is a msgpack message
		for _, part := range parts {
			var dec *codec.Decoder = codec.NewDecoderBytes(part, &h)
			err := dec.Decode(&msg)
			checkError(err)
			unpacked := unpack.HandshakeEvt(msg)
			fmt.Println(&unpacked)
		}
	} else {
		fmt.Println("received ", string(joinedBytes))
	}
}
