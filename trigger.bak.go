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
		unpackEvent(msg, part)
		//unpacked.Print()
	}
}

func convertMapToIntKeys(event map[interface{}]interface{}) map[int64]interface{} {
	anotherEvt := make(map[int64]interface{})

	for key, value := range event {
		switch value.(type) {
		case map[interface{}]interface{}:
			anotherEvt[key.(int64)] = interface{}(convertMapToIntKeys(value.(map[interface{}]interface{})))
		default:
			anotherEvt[key.(int64)] = value
		}
	}

	return anotherEvt
}

func parseEventType(event interface{}) int64 {
	var evtType int64

	for k, v := range event.(map[interface{}]interface{}) {
		if k.(int64) == 0 {
			evtType = v.(int64)
		}
	}

	return evtType
}

type ObjCreated struct {
	EventType int64
	Timestamp float64
	Payload   objCreatedPayload
}

type objCreatedPayload struct {
	ObjectId  uint64
	ClassName string
}

func ObjCreatedEvt(payload map[int64]interface{}) (response ObjCreated) {
	payloadMap := payload[2].(map[int64]interface{})

	objCreatedPayloadObj := objCreatedPayload{
		ObjectId:  payloadMap[3].(uint64),
		ClassName: payloadMap[4].(string),
	}

	response = ObjCreated{
		EventType: payload[0].(int64),
		Timestamp: payload[1].(float64),
		Payload:   objCreatedPayloadObj,
	}

	return
}

func unpackEvent(payload map[int]interface{}, part []byte) {
	eventType := payload[0].(int64)
	/* This event is always of the type Event Collection */

	if eventType != 7 {
		panic("Event other than event collection has been received. Was the protocol changed? Aborting")
	}

	evtCollectionPayload := payload[2].([]interface{})

	for _, event := range evtCollectionPayload {
		/* get event type. because we can't rely on int keys for events
		   because gc stats doesn't have int keys inside payload */
		evtType := parseEventType(event)

		switch evtType {
		case 0:
			evt := convertMapToIntKeys(event.(map[interface{}]interface{}))
			unpackedEvent := ObjCreatedEvt(evt)
			fmt.Println(unpackedEvent)
		}
		//fmt.Println(unpackedEvent)

		//unpackEvent.Print()
		//case 1:
		//unpackedEvent = unpack.ObjDestroyedEvt(event)
		//case 2:
		//unpackedEvent = unpack.GcStartEvt(event)
		//case 3:
		//unpackedEvent = unpack.GcEndMinorEvt(event)
		//case 4:
		//unpackedEvent = unpack.GcEndSweepEvt(event)
		////case 5:
		////[> Not implemented <]
		////unpack.ObjectSpaceDumpEvt(event)
		//case 6:
		//unpackedEvent = unpack.GcStatsEvt(event)
		//}
	}

	/* iterate over the payloadMap object and check the event of each
	element. This would enable us to properly unpack and generate objects */
	//for _, event := range evtCollectionPayload {
	//var unpackedEvent EventInterface
	//fmt.Println(event)
	//eventType := event[0]

	//append(events, unpackedEvent)
	//}
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
