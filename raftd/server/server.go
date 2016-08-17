package server

import (
	"github.com/xsdb/playground/raftd/partition"
	"github.com/xsdb/playground/raftd/xsmember"
)

type Server struct {
	p        *partition.Partition
	xsmember *xsmember.Xsmember
	eventCh  chan *xsmember.Event
}

func NewServer() (*Server, error) {
	s := &Server{}

	s.setupXsmember()

	s.p, _ = partition.NewPartition()

	return s, nil
}

func (s *Server) setupXsmember() {
	s.eventCh = make(chan *xsmember.Event, 256)

	conf := &xsmember.Config{eventCh: s.eventCh}

	xsmember, err := xsmember.NewXsmember(s.conf)
	if err != nil {
		log.Printf("xsmember create fail. %v", err)
	}

	s.xsmember = xsmember

	go s.xsMemberEventLoop()
}

func (s *Server) xsMemberEventLoop() {
	for {
		select {
		case e := <-s.eventCh:
			s.handleMemberEvent(e)
		}
	}
}

func (s *Server) handleMemberEvent(e *xsmember.Event) {
	switch e.Type() {
	case xsmember.TypeJoin:
	case xsmember.TypeLeave:
	case xsmember.TypeUpdate:
	}
}

func (s *Server) Start() {
}
