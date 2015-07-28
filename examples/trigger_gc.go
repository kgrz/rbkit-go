/* This program simply triggers a GC every 1 second and reads and parses the
* response. Useful for quick debugging of the network communication.
 */
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
		case parts := <-dataChans.In():
			go func() {
				fmt.Println("Got info")
				processIncomingEvent(parts)
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
	var msg map[int]interface{}
	joinedBytes := bytes.Join(parts, nil)

	if len(joinedBytes) != 2 {
		// This is a msgpack message
		for _, part := range parts {
			var dec *codec.Decoder = codec.NewDecoderBytes(part, &h)
			err := dec.Decode(&msg)
			checkError(err)
			unpacked := unpackHandshake(msg)
			fmt.Println(&unpacked)
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

func (h *HandshakeResponse) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Event Type: Handshake\n")
	buffer.WriteString(fmt.Sprintln("Timestamp : ", time.Unix(int64(h.Timestamp), 0)))
	buffer.WriteString(fmt.Sprintln("Rbkit Server Version: ", h.Payload.ServerVersion))
	buffer.WriteString(fmt.Sprintln("Rbkit Protocol Version: ", h.Payload.ProtocolVersion))
	buffer.WriteString(fmt.Sprintln("Process Name: ", h.Payload.ProcessName))
	buffer.WriteString(fmt.Sprintln("Working Directory: ", h.Payload.Pwd))
	buffer.WriteString(fmt.Sprintln("Pid: ", h.Payload.Pid))

	if h.Payload.ObjectTraceEnabled {
		buffer.WriteString("Object Trace Enabled\n")
	} else {
		buffer.WriteString("Object Trace Not Enabled\n")
	}
	return buffer.String()
}
