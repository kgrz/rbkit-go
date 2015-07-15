package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/kgrz/rbkit-go/unpack"
	"github.com/ugorji/go/codec"
	"github.com/vaughan0/go-zmq"
)

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
			unpacked.Print()
		}
	} else {
		fmt.Println("received ", string(joinedBytes))
	}
}

/* This handles processing event data */
func processIncomingEvent(parts [][]byte) {
	var msg map[int]interface{}

	for _, part := range parts {
		var dec *codec.Decoder = codec.NewDecoderBytes(part, &h)
		err := dec.Decode(&msg)
		checkError(err)
		unpacked := unpackEvent(msg, part)
		unpacked.Print()
	}
}

func unpackEvent(payload map[int]interface{}, part []byte) (events []EventInterface) {
	eventType := payload[0]
	payloadData := payload[2]
	/* This event is always of the type Event Collection */

	if eventType != 7 {
		panic("Event other than event collection has been received. Was the protocol changed? Aborting")
	}

	evtCollectionPayload := payload[2].(map[interface{}]interface{})

	/* iterate over the payloadMap object and check the event of each
	element. This would enable us to properly unpack and generate objects */
	for _, event := range evtCollectionPayload {
		var unpackedEvent EventInterface
		fmt.Println(event)
		eventType := event[0]

		switch eventType {
		case 0:
			unpackedEvent = unpack.ObjCreatedEvt(event)
		case 1:
			unpackedEvent = unpack.ObjDestroyedEvt(event)
		case 2:
			unpackedEvent = unpack.GcStartEvt(event)
		case 3:
			unpackedEvent = unpack.GcEndMinorEvt(event)
		case 4:
			unpackedEvent = unpack.GcEndSweepEvt(event)
		//case 5:
		//[> Not implemented <]
		//unpack.ObjectSpaceDumpEvt(event)
		case 6:
			unpackedEvent = unpack.GcStatsEvt(event)
		}

		append(events, unpackedEvent)
	}
	/* handshake is not from sub socket. this has to be handled inside
	* command channel */
	//case 8:
	//unpack.HandshakeEvt(payload)
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

type EventInterface interface {
	Print()
}
