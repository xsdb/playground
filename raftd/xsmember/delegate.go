package xsmember

import (
	"fmt"
	"log"
)

type delegate struct {
	xsmember *Xsmember
}

func (d *delegate) NodeMeta(limit int) []byte {
	bytes := d.xsmember.encodeTags(d.xsmember.conf.Tags)
	if len(bytes) > limit {
		panic(fmt.Errorf("Node tags '%v' exceeds length limit of %d bytes", d.xsmember.conf.Tags, limit))
	}

	return bytes
}

func (d *delegate) NotifyMsg(in []byte) {
	if len(in) == 0 {
		return
	}

	log.Printf("NotifyMsg %v", in)

	t := msgType(in[0])
	switch t {
	case MsgTypeAddPeer:
		addPeerMsg := &MsgAddPeer{}
		decodeMessage(in[1:], addPeerMsg)

		log.Printf("AddPeer %v, %v", addPeerMsg.Addr, addPeerMsg.Port)

		msg := &Message{Type: t, Body: addPeerMsg}
		d.xsmember.conf.msgCh <- msg
	}
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	msgs := d.xsmember.broadcasts.GetBroadcasts(overhead, limit)

	return msgs
}

func (d *delegate) LocalState(join bool) []byte {
	buf := make([]byte, 1)

	return buf
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {
}
