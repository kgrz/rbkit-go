// Prompt user for commands, print them out. This has to be done in a
// separate thread
package main

import (
	"bytes"
	"fmt"
	"time"

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
	var h codec.MsgpackHandle
	h.WriteExt = true
	h.RawToString = true
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
		var msg map[int]interface{}
		optionChan := askForOption()
		chans.Out() <- [][]byte{[]byte(<-optionChan)}

		select {
		case parts := <-chans.In():
			joinedBytes := bytes.Join(parts, nil)

			if len(joinedBytes) != 2 {
				// This is a msgpack message
				for _, part := range parts {
					var dec *codec.Decoder = codec.NewDecoderBytes(part, &h)
					err = dec.Decode(&msg)
					checkError(err)
					unpacked := unpackHandshake(msg)
					unpacked.Print()
				}
			} else {
				fmt.Println("received ", string(joinedBytes))
			}
		case err := <-chans.Errors():
			checkError(err)
		}
	}
}

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

func parseReponse(parts [][]byte, h codec.MsgpackHandle, msg map[int]interface{}) {
	joinedBytes := bytes.Join(parts, nil)

	if len(joinedBytes) != 2 {
		// This is a msgpack message
		for _, part := range parts {
			var dec *codec.Decoder = codec.NewDecoderBytes(part, &h)
			err := dec.Decode(&msg)
			checkError(err)
			unpacked := unpackHandshake(msg)
			unpacked.Print()
		}
	} else {
		fmt.Println("\n")
		fmt.Println("received ", string(joinedBytes))
		fmt.Println("\n")
	}

}

func unpackHandshake(payload map[int]interface{}) (response HandshakeResponse) {
	payloadMap := payload[2].(map[interface{}]interface{})
	var objectTraceEnabled bool

	if payloadMap["object_trace_enabled"].(int64) == 0 {
		objectTraceEnabled = false
	} else {
		objectTraceEnabled = true
	}

	handshakePayloadObj := handshakePayload{
		ServerVersion:      payloadMap["rbkit_server_version"].(string),
		ProtocolVersion:    payloadMap["rbkit_protocol_version"].(string),
		ProcessName:        payloadMap["process_name"].(string),
		Pwd:                payloadMap["pwd"].(string),
		Pid:                payloadMap["pid"].(uint64),
		ObjectTraceEnabled: objectTraceEnabled,
	}

	response = HandshakeResponse{
		EventType: payload[0].(int64),
		Timestamp: payload[1].(float64),
		Payload:   handshakePayloadObj,
	}

	return
}

func (h HandshakeResponse) Print() {
	fmt.Println("\n")
	fmt.Println("Event Type: Handshake")
	fmt.Println("Timestamp : ", time.Unix(int64(h.Timestamp), 0))
	fmt.Println("Rbkit Server Version: ", h.Payload.ServerVersion)
	fmt.Println("Rbkit Protocol Version: ", h.Payload.ProtocolVersion)
	fmt.Println("Process Name: ", h.Payload.ProcessName)
	fmt.Println("Working Directory: ", h.Payload.Pwd)
	fmt.Println("Pid: ", h.Payload.Pid)

	if h.Payload.ObjectTraceEnabled {
		fmt.Println("Object Trace Enabled")
	} else {
		fmt.Println("Object Trace Not Enabled")
	}
	fmt.Println("\n")
}
