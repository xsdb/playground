package server

import (
	"github.com/hashicorp/memberlist"
	"github.com/xsdb/playground/raftd/partition"
)

type Server struct {
	p    *partition.Partition
	list *memberlist.Memberlist
}

func NewServer() (*Server, error) {
	s := &Server{}

	s.list = s.setupMemberlist()

	s.p, _ = partition.NewPartition()

	return s, nil
}

func (s *Server) setupMemberlist() *memberlist.Memberlist {
	conf := memberlist.DefaultLocalConfig()

	conf.Events = &eventDelegate{server: s}

	list, err := memberlist.Create(conf)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	return list
}

func (s *Server) handleNodeJoin(n *memberlist.Node) {
}

func (s *Server) handleNodeLeave(n *memberlist.Node) {
}

func (s *Server) handleNodeUpdate(n *memberlist.Node) {
}

func (s *Server) Start() {
}
