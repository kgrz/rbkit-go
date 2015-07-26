package unpack

import (
	"bytes"
	"fmt"
	"reflect"
	"time"

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

type ObjDestroyed struct {
	EventType int64
	Timestamp float64
	Payload   objDestroyedPayload
}

type objDestroyedPayload struct {
	ObjectId uint64
}

func ObjDestroyedEvt(payload map[int64]interface{}) (response ObjDestroyed) {
	payloadMap := payload[2].(map[int64]interface{})

	objDestroyedPayloadObj := objDestroyedPayload{
		ObjectId: payloadMap[3].(uint64),
	}

	response = ObjDestroyed{
		EventType: payload[0].(int64),
		Timestamp: payload[1].(float64),
		Payload:   objDestroyedPayloadObj,
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

func GcStartEvt(payload map[int64]interface{}) (response GcStart) {
	response = GcStart{
		EventType: payload[0].(int64),
		Timestamp: payload[1].(float64),
	}

	return
}

func GcEndSweepEvt(payload map[int64]interface{}) (response GcEndSweep) {
	response = GcEndSweep{
		EventType: payload[0].(int64),
		Timestamp: payload[1].(float64),
	}

	return
}

func GcEndMinorEvt(payload map[int64]interface{}) (response GcEndMinor) {
	response = GcEndMinor{
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

func GcStatsEvt(event map[int64]interface{}) (response GcStat) {
	/* we try to encode the payload part first to msgpack format and then
	* decode it again so that we get the struct conversion easily instead of
	* manually assigning each element of the struct to the payload's value item
	 *
	 *  This is foolish
	 *
	 * The reason for these gymnastics is because the keys of the normal
	 * message are ints whereas the keys of gc stats are strings. Go msgpack
	 * can't parse int key objects if it's given a struct to parse to. So we
	 * take out the payload portion, encode it, and then decode it back to the
	 * struct format in one shot by giving the target type.
	*/
	payload := event[2].(map[string]interface{})

	var b []byte = make([]byte, 0, 100)
	var h codec.MsgpackHandle
	h.WriteExt = true
	h.RawToString = true
	h.MapType = reflect.TypeOf(map[string]interface{}(nil))
	var enc *codec.Encoder = codec.NewEncoderBytes(&b, &h)
	err := enc.Encode(payload)

	var dec *codec.Decoder = codec.NewDecoderBytes(b, &h)
	err = dec.Decode(&response)
	checkError(err)

	return
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (evt *GcEndMinor) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Event Type : GC End Minor")
	buffer.WriteString(fmt.Sprintln("Timestamp : ", time.Unix(int64(evt.Timestamp), 0)))
	return buffer.String()
}

func (evt *GcEndSweep) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Event Type : GC End Sweep")
	buffer.WriteString(fmt.Sprintln("Timestamp : ", time.Unix(int64(evt.Timestamp), 0)))
	return buffer.String()
}

func (evt *GcStart) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Event Type : GC Start")
	buffer.WriteString(fmt.Sprintln("Timestamp : ", time.Unix(int64(evt.Timestamp), 0)))
	return buffer.String()
}

func (evt *GcStat) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Event Type : GC Start")
	buffer.WriteString(fmt.Sprintln("GC Count", evt.Count))
	buffer.WriteString(fmt.Sprintln("Minor GC Count", evt.MinorGcCount))
	buffer.WriteString(fmt.Sprintln("Major GC Count", evt.MajorGcCount))
	buffer.WriteString(fmt.Sprintln("Heap Allocated Pages", evt.HeapAllocatedPages))
	buffer.WriteString(fmt.Sprintln("Heap Eden Pages", evt.HeapEdenPages))
	buffer.WriteString(fmt.Sprintln("Heap Tomb Pages", evt.HeapTombPages))
	buffer.WriteString(fmt.Sprintln("Heap Allocatable Pages", evt.HeapAllocatablePages))
	buffer.WriteString(fmt.Sprintln("Heap Sorted Length", evt.HeapSortedLength))
	buffer.WriteString(fmt.Sprintln("Heap Live Slots", evt.HeapLiveSlots))
	buffer.WriteString(fmt.Sprintln("Heap Free Slots", evt.HeapFreeSlots))
	buffer.WriteString(fmt.Sprintln("Heap Sweep Slots", evt.HeapSweptSlots))
	buffer.WriteString(fmt.Sprintln("Old Objects", evt.OldObjects))
	buffer.WriteString(fmt.Sprintln("Old Objects Limit", evt.OldObjectsLimit))
	buffer.WriteString(fmt.Sprintln("Total Allocated Objects", evt.TotalAllocatedObjects))
	buffer.WriteString(fmt.Sprintln("Total Freed Objects", evt.TotalFreedObjects))
	buffer.WriteString(fmt.Sprintln("Heap Final Slots", evt.HeapFinalSlots))
	buffer.WriteString(fmt.Sprintln("Malloc Increased Bytes", evt.MallocIncreaseBytes))
	buffer.WriteString(fmt.Sprintln("Malloc Increased Bytes Limit", evt.MallocIncreaseBytesLimit))
	buffer.WriteString(fmt.Sprintln("Old Malloc Increased Bytes", evt.OldMallocIncreaseBytes))
	buffer.WriteString(fmt.Sprintln("Old Malloc Increased Bytes Limit", evt.OldMallocIncreaseBytesLimit))
	buffer.WriteString(fmt.Sprintln("Total Heap Size", evt.TotalHeapSize))
	buffer.WriteString(fmt.Sprintln("Total Memsize", evt.TotalMemsize))
	return buffer.String()
}

func (evt *ObjCreated) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Event Type : Object Created")
	buffer.WriteString(fmt.Sprintln("Timestamp : ", time.Unix(int64(evt.Timestamp), 0)))
	buffer.WriteString(fmt.Sprintln("Object ID: ", evt.Payload.ObjectId))
	buffer.WriteString(fmt.Sprintln("Class Name: ", evt.Payload.ClassName))
	return buffer.String()
}

func (evt *ObjDestroyed) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Event Type : Object Destroyed")
	buffer.WriteString(fmt.Sprintln("Timestamp : ", time.Unix(int64(evt.Timestamp), 0)))
	buffer.WriteString(fmt.Sprintln("Object ID: ", evt.Payload.ObjectId))
	return buffer.String()
}

func (h *Handshake) String() string {
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
