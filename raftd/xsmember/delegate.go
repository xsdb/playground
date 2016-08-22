package xsmember

import (
	"fmt"
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

func (d *delegate) NotifyMsg([]byte) {
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	buf := make([][]byte, 1)

	return buf
}

func (d *delegate) LocalState(join bool) []byte {
	buf := make([]byte, 1)

	return buf
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {
}
