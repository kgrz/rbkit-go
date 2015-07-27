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

	commandChans := commandSock.Channels()
	defer commandChans.Close()

	dataChans := dataSock.Channels()
	defer dataChans.Close()

	// send the command once. Do this inside a go routine
	//sendCommand(chans)

	for {
		select {
		case parts := <-commandChans.In():
			go func() {
				processIncomingMessage(parts)
				//sendCommand(chans)
			}()
		case err := <-commandChans.Errors():
			checkError(err)
		case parts := <-dataChans.In():
			go func() {
				processIncomingEvent(parts)
			}()
		case err := <-dataChans.Errors():
			checkError(err)
		}
		//time.Sleep(time.Second * 1)
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
			unpacked := unpack.HandshakeEvt(msg)
			fmt.Println(unpacked)
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
		unpackEventCollection(msg, part)
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

func unpackEventCollection(payload map[int]interface{}, part []byte) {
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

		var unpacked Event

		switch evtType {
		case 0:
			evt := convertMapToIntKeys(event.(map[interface{}]interface{}))
			unpackedEvent := unpack.ObjCreatedEvt(evt)
			unpacked = &unpackedEvent
		case 1:
			evt := convertMapToIntKeys(event.(map[interface{}]interface{}))
			unpackedEvent := unpack.ObjDestroyedEvt(evt)
			unpacked = &unpackedEvent
		case 2:
			evt := convertMapToIntKeys(event.(map[interface{}]interface{}))
			unpackedEvent := unpack.GcStartEvt(evt)
			unpacked = &unpackedEvent
		case 3:
			evt := convertMapToIntKeys(event.(map[interface{}]interface{}))
			unpackedEvent := unpack.GcEndMinorEvt(evt)
			unpacked = &unpackedEvent
		case 4:
			evt := convertMapToIntKeys(event.(map[interface{}]interface{}))
			unpackedEvent := unpack.GcEndSweepEvt(evt)
			unpacked = &unpackedEvent
		case 6:
			evt := convertGcStatsEvtKeys(event.(map[interface{}]interface{}))
			unpackedEvent := unpack.GcStatsEvt(evt)
			unpacked = &unpackedEvent
		}

		fmt.Println(unpacked)
	}
}

type Event interface {
	String() string
}
