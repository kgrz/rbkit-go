package unpack

import (
	"fmt"

	"github.com/ugorji/go/codec"
)

type Handshake struct {
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

func HandshakeEvt(payload map[int]interface{}) (response Handshake) {
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

	response = Handshake{
		EventType: payload[0].(int64),
		Timestamp: payload[1].(float64),
		Payload:   handshakePayloadObj,
	}

	return
}

type ObjCreated struct {
	EventType int64
	Timestamp float64
	Payload   objCreatedPayload
}

type objCreatedPayload struct {
	ObjectId  int64
	ClassName string
}

func ObjCreatedEvt(payload map[int]interface{}) (response ObjCreated) {
	payloadMap := payload[2].(map[interface{}]interface{})

	objCreatedPayloadObj := objCreatedPayload{
		ObjectId:  payloadMap[3].(int64),
		ClassName: payloadMap[3].(string),
	}

	response = ObjCreated{
		EventType: payload[0].(int64),
		Timestamp: payload[1].(float64),
		Payload:   objCreatedPayloadObj,
	}

	return
}

type ObjDestroyed struct {
	EventType int64
	Timestamp float64
	Payload   objDestroyedPayload
}

type objDestroyedPayload struct {
	ObjectId int64
}

func ObjDestroyedEvt(payload map[int]interface{}) (response ObjDestroyed) {
	payloadMap := payload[2].(map[interface{}]interface{})

	objCreatedPayloadObj := objCreatedPayload{
		ObjectId: payloadMap[3].(int64),
	}

	reponse := ObjCreated{
		EventType: payload[0].(int64),
		Timestamp: payload[1].(float64),
		Payload:   objCreatedPayloadObj,
	}

	return
}

type GcStart struct {
	EventType int64
	Timestamp float64
}

type GcEndSweep struct {
	EventType int64
	Timestamp float64
}

type GcEndMinor struct {
	EventType int64
	Timestamp float64
}

func GcStartEvt(payload interface{}) (response GcStart) {
	payloadMap := payload[2].(map[interface{}]interface{})

	reponse := GcStart{
		EventType: payload[0].(int64),
		Timestamp: payload[1].(float64),
	}

	return
}

func GcEndSweepEvt(payload map[int]interface{}) (response GcEndSweep) {
	payloadMap := payload[2].(map[interface{}]interface{})

	reponse := GcEndSweep{
		EventType: payload[0].(int64),
		Timestamp: payload[1].(float64),
	}

	return
}

func GcEndMinorEvt(payload map[int]interface{}) (response GcEndMinor) {
	payloadMap := payload[2].(map[interface{}]interface{})

	reponse := GcEndMinor{
		EventType: payload[0].(int64),
		Timestamp: payload[1].(float64),
	}

	return
}

type GcStat struct {
	Count                       int64 `codec:"count"`
	MinorGcCount                int64 `codec:"minor_gc_count"`
	MajorGcCount                int64 `codec:"major_gc_count"`
	HeapAllocatedPages          int64 `codec:"heap_allocated_pages"`
	HeapEdenPages               int64 `codec:"heap_eden_pages"`
	HeapTombPages               int64 `codec:"heap_tomb_pages"`
	HeapAllocatablePages        int64 `codec:"heap_allocatable_pages"`
	HeapSortedLength            int64 `codec:"heap_sorted_length"`
	HeapLiveSlots               int64 `codec:"heap_live_slots"`
	HeapFreeSlots               int64 `codec:"heap_free_slots"`
	HeapSweptSlots              int64 `codec:"heap_swept_slots"`
	OldObjects                  int64 `codec:"old_objects"`
	OldObjectsLimit             int64 `codec:"old_objects_limit"`
	TotalAllocatedObjects       int64 `codec:"total_allocated_objects"`
	TotalFreedObjects           int64 `codec:"total_freed_objects"`
	HeapFinalSlots              int64 `codec:"heap_final_slots"`
	MallocIncreaseBytes         int64 `codec:"malloc_increase_bytes"`
	MallocIncreaseBytesLimit    int64 `codec:"malloc_increase_bytes_limit"`
	OldMallocIncreaseBytes      int64 `codec:"oldmalloc_increase_bytes"`
	OldMallocIncreaseBytesLimit int64 `codec:"oldmalloc_increase_bytes_limit"`
	TotalHeapSize               int64 `codec:"total_heap_size"`
	TotalMemsize                int64 `codec:"total_memsize"`
}

func GcStatsEvt(payload map[int]interface{}) (response GcStat) {
	var h codec.MsgpackHandle
	h.WriteExt = true
	h.RawToString = true
	part := payload[2].([]byte)
	fmt.Println(part)

	var dec *codec.Decoder = codec.NewDecoderBytes(part, &h)
	err := dec.Decode(&response)
	checkError(err)

	return
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (evt GcEndMinor) Print() {
	fmt.Println("GC End minor")
}

func (evt GcEndSweep) Print() {
	fmt.Println("GC End sweep")
}

func (evt GcStart) Print() {
	fmt.Println("GC Start")
}

func (evt GcStat) Print() {
	fmt.Println("GC Stat")
}

func (evt ObjCreated) Print() {
	fmt.Println("Object Created")
}

func (evt ObjDestroyed) Print() {
	fmt.Println("Object Destroyed")
}
