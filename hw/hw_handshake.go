package hw

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/kgrz/msgp/msgp"
)

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Handshake) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var intField int64
	_ = intField
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		intField, bts, err = msgp.ReadInt64MapKeyZC(bts)
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

// UnmarshalMsg implements msgp.Unmarshaler
func (z *handshakePayload) UnmarshalMsg(bts []byte) (o []byte, err error) {
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

func (z *handshakePayload) Msgsize() (s int) {
	s = 1 + 21 + msgp.StringPrefixSize + len(z.ServerVersion) + 23 + msgp.StringPrefixSize + len(z.ProtocolVersion) + 13 + msgp.StringPrefixSize + len(z.ProcessName) + 4 + msgp.StringPrefixSize + len(z.Pwd) + 4 + msgp.Uint64Size + 21 + msgp.BoolSize + 22 + msgp.BoolSize + 11 + msgp.StringPrefixSize + len(z.ClockType) + 19 + msgp.StringPrefixSize + len(z.CpuProfilngMode)
	return
}
