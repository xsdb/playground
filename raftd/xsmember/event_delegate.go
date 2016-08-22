package xsmember

import (
	"log"
	"net"

	"github.com/hashicorp/memberlist"
)

type eventDelegate struct {
	xsmember *Xsmember
}

func (e *eventDelegate) NotifyJoin(n *memberlist.Node) {
	log.Printf("Notify Join %v", n)

	m := e.nodeToMember(n)
	e.xsmember.handleMemberJoin(m)
}

func (e *eventDelegate) NotifyLeave(n *memberlist.Node) {
	log.Printf("Notify Leave %v", n)

	m := e.nodeToMember(n)
	e.xsmember.handleMemberLeave(m)
}

func (e *eventDelegate) NotifyUpdate(n *memberlist.Node) {
	log.Printf("Notify Update %v", n)

	m := e.nodeToMember(n)
	e.xsmember.handleMemberUpdate(m)
}

func (e *eventDelegate) nodeToMember(n *memberlist.Node) *Member {
	return &Member{
		Name:        n.Name,
		Addr:        net.IP(n.Addr),
		Port:        n.Port,
		Tags:        e.xsmember.decodeTags(n.Meta),
		Status:      StatusNone,
		ProtocolMin: n.PMin,
		ProtocolMax: n.PMax,
		ProtocolCur: n.PCur,
		DelegateMin: n.DMin,
		DelegateMax: n.DMax,
		DelegateCur: n.DCur,
	}
}
