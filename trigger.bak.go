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
				processIncomingMessage(parts)
				sendCommand(chans)
			}()
		case err := <-chans.Errors():
			checkError(err)
		case parts := <-dataChans.In():
			go func() {
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
		chans.Out() <- [][]byte{[]byte("trigger_gc")}
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
		//fmt.Println("received ", string(joinedBytes))
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

func convertGcStatsEvtKeys(event map[interface{}]interface{}) map[int64]interface{} {
	gcStats := make(map[int64]interface{})

	for key, value := range event {
		switch value.(type) {
		case map[interface{}]interface{}:
			gcStats[key.(int64)] = interface{}(convertMapToStringKeys(value.(map[interface{}]interface{})))
		default:
			gcStats[key.(int64)] = value
		}
	}

	return gcStats
}

func convertMapToStringKeys(value map[interface{}]interface{}) map[string]interface{} {
	gcStats := make(map[string]interface{})

	for key, value := range value {
		switch value.(type) {
		case map[interface{}]interface{}:
			gcStats[key.(string)] = interface{}(convertMapToStringKeys(value.(map[interface{}]interface{})))
		default:
			gcStats[key.(string)] = value
		}
	}

	return gcStats
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
			unpack.ObjCreatedEvt(evt)
		case 1:
			evt := convertMapToIntKeys(event.(map[interface{}]interface{}))
			unpack.ObjDestroyedEvt(evt)
		case 2:
			evt := convertMapToIntKeys(event.(map[interface{}]interface{}))
			unpackedEvent := unpack.GcStartEvt(evt)
			fmt.Println("gc started")
			fmt.Println(unpackedEvent)
		case 3:
			evt := convertMapToIntKeys(event.(map[interface{}]interface{}))
			unpackedEvent := unpack.GcEndMinorEvt(evt)
			fmt.Println("gc end minor")
			fmt.Println(unpackedEvent)
		case 4:
			evt := convertMapToIntKeys(event.(map[interface{}]interface{}))
			unpackedEvent := unpack.GcEndSweepEvt(evt)
			fmt.Println("gc end sweep")
			fmt.Println(unpackedEvent)
		case 6:
			fmt.Println("gc stats")
			evt := convertGcStatsEvtKeys(event.(map[interface{}]interface{}))
			unpackedEvent := unpack.GcStatsEvt(evt)
			fmt.Println(unpackedEvent)
		}
		//unpackEvent.Print()
	}
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
