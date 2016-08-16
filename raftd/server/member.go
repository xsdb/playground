package server

import (
	"github.com/hashicorp/memberlist"
)

type eventDelegate struct {
	server *Server
}

func (e *eventDelegate) NotifyJoin(n *memberlist.Node) {
	e.server.handleNodeJoin(n)
}

func (e *eventDelegate) NotifyLeave(n *memberlist.Node) {
	e.server.handleNodeLeave(n)
}

func (e *eventDelegate) NotifyUpdate(n *memberlist.Node) {
	e.server.handleNodeUpdate(n)
}
