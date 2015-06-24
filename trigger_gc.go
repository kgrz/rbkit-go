package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ugorji/go/codec"
	"github.com/vaughan0/go-zmq"
)

type HandshakeResponse struct {
	EventType int64
	Timestamp float64
	Payload   kandshakePayload
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
	var msg map[int]interface{}

	ctx, err := zmq.NewContext()
	checkError(err)
	defer ctx.Close()

	commandSock, err := ctx.Socket(zmq.Req)
	checkError(err)
	defer commandSock.Close()

	err = commandSock.Connect("tcp://127.0.0.1:5556")
	checkError(err)

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
				unpacked := unpackHandshake(msg)
				fmt.Println(unpacked)
				//printHandshakeInfo(msg)
			}
		} else {
			fmt.Println("received ", string(joinedBytes))
		}

		time.Sleep(time.Second * 1)
	}
}

func printHandshakeInfo(payload map[int]interface{}) {
	// unpack the interface and convert the value to an int
	timestamp := payload[1].(int64)
	// unpack the payload interface and convert the value to map that can then
	// be queried using the string keys for respective fields
	payloadMap := payload[2].(map[interface{}]interface{})

	fmt.Println("\n")
	fmt.Println("Event Type: Handshake")
	fmt.Println("Timestamp : ", time.Unix(int64(timestamp), 0))
	fmt.Println("Rbkit Server Version: ", payloadMap["rbkit_server_version"])
	fmt.Println("Rbkit Protocol Version: ", payloadMap["rbkit_protocol_version"])
	fmt.Println("Process Name: ", payloadMap["process_name"])
	fmt.Println("Working Directory: ", payloadMap["pwd"])
	fmt.Println("Pid: ", payloadMap["pid"])

	if payloadMap["object_trace_enabled"] == 0 {
		fmt.Println("Object Trace Not Enabled")
	} else {
		fmt.Println("Object Trace Enabled")
	}
	fmt.Println("\n")
}

func unpackHandshake(payload map[int]interface{}) (response HandshakeResponse) {
	payloadMap := payload[2].(map[interface{}]interface{})
	var objectTraceEnabled bool

	if payloadMap["object_trace_enabled"] == 0 {
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
