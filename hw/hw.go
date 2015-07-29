package hw

//go:generate msgp -io=false -o hw_gen.go

type Handshake struct {
	EventType int64
	Timestamp float64
	Payload   HandshakePayload
}

type HandshakePayload struct {
	ServerVersion       string `msg:"rbkit_server_version"`
	ProtocolVersion     string `msg:"rbkit_protocol_version"`
	ProcessName         string `msg:"process_name"`
	Pwd                 string `msg:"pwd"`
	Pid                 uint64 `msg:"pid"`
	ObjectTraceEnabled  int64  `msg:"object_trace_enabled"`
	CpuProfilingEnabled int64  `msg:"cpu_profiling_enabled"`
	ClockType           string `msg:"clock_type"`
	CpuProfilngMode     string `msg:"cpu_profiling_mode"`
}
