package xsmember

import (
	"bytes"
	"github.com/hashicorp/go-msgpack/codec"
)

type MsgAddPeer struct {
	Addr string
	Port int
}

type Message struct {
	Type msgType
	Body interface{}
}

type msgType uint8

const (
	MsgTypeNone msgType = iota
	MsgTypeAddPeer
)

func decodeMessage(buf []byte, out interface{}) error {
	var handle codec.MsgpackHandle
	return codec.NewDecoder(bytes.NewReader(buf), &handle).Decode(out)
}

func encodeMessage(t msgType, msg interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(uint8(t))

	handle := codec.MsgpackHandle{}
	encoder := codec.NewEncoder(buf, &handle)
	err := encoder.Encode(msg)
	return buf.Bytes(), err
}
