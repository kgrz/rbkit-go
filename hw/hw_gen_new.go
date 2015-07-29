package hw

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z *Handshake) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "EventType"
	o = append(o, 0x83, 0xa9, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65)
	o = msgp.AppendInt64(o, z.EventType)
	// string "Timestamp"
	o = append(o, 0xa9, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70)
	o = msgp.AppendFloat64(o, z.Timestamp)
	// string "Payload"
	o = append(o, 0xa7, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64)
	o, err = z.Payload.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Handshake) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	var intField int64
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		intField, bts, err = msgp.ReadMapInt64KeyZC(bts)
		if err != nil {
			return
		}
		switch intField {
		case 0:
			z.EventType, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		case 1:
			z.Timestamp, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case 2:
			bts, err = z.Payload.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *Handshake) Msgsize() (s int) {
	s = 1 + 10 + msgp.Int64Size + 10 + msgp.Float64Size + 8 + z.Payload.Msgsize()
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *HandshakePayload) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 9
	// string "rbkit_server_version"
	o = append(o, 0x89, 0xb4, 0x72, 0x62, 0x6b, 0x69, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.ServerVersion)
	// string "rbkit_protocol_version"
	o = append(o, 0xb6, 0x72, 0x62, 0x6b, 0x69, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.ProtocolVersion)
	// string "process_name"
	o = append(o, 0xac, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.ProcessName)
	// string "pwd"
	o = append(o, 0xa3, 0x70, 0x77, 0x64)
	o = msgp.AppendString(o, z.Pwd)
	// string "pid"
	o = append(o, 0xa3, 0x70, 0x69, 0x64)
	o = msgp.AppendUint64(o, z.Pid)
	// string "object_trace_enabled"
	o = append(o, 0xb4, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64)
	o = msgp.AppendInt64(o, z.ObjectTraceEnabled)
	// string "cpu_profiling_enabled"
	o = append(o, 0xb5, 0x63, 0x70, 0x75, 0x5f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64)
	o = msgp.AppendInt64(o, z.CpuProfilingEnabled)
	// string "clock_type"
	o = append(o, 0xaa, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendString(o, z.ClockType)
	// string "cpu_profiling_mode"
	o = append(o, 0xb2, 0x63, 0x70, 0x75, 0x5f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x6f, 0x64, 0x65)
	o = msgp.AppendString(o, z.CpuProfilngMode)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *HandshakePayload) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "rbkit_server_version":
			z.ServerVersion, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "rbkit_protocol_version":
			z.ProtocolVersion, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "process_name":
			z.ProcessName, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "pwd":
			z.Pwd, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "pid":
			z.Pid, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "object_trace_enabled":
			z.ObjectTraceEnabled, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		case "cpu_profiling_enabled":
			z.CpuProfilingEnabled, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		case "clock_type":
			z.ClockType, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "cpu_profiling_mode":
			z.CpuProfilngMode, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *HandshakePayload) Msgsize() (s int) {
	s = 1 + 21 + msgp.StringPrefixSize + len(z.ServerVersion) + 23 + msgp.StringPrefixSize + len(z.ProtocolVersion) + 13 + msgp.StringPrefixSize + len(z.ProcessName) + 4 + msgp.StringPrefixSize + len(z.Pwd) + 4 + msgp.Uint64Size + 21 + msgp.BoolSize + 22 + msgp.BoolSize + 11 + msgp.StringPrefixSize + len(z.ClockType) + 19 + msgp.StringPrefixSize + len(z.CpuProfilngMode)
	return
}
