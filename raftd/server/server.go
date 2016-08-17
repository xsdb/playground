package server

import (
	"log"

	"github.com/xsdb/playground/raftd/partition"
	"github.com/xsdb/playground/raftd/partitionManager"
	"github.com/xsdb/playground/raftd/xsmember"
)

type Server struct {
	p *partition.Partition
	/* for cluster member manage */
	xsmember *xsmember.Xsmember

	/* for recv cluster member event */
	eventCh chan *xsmember.Event

	/* for partition manage
	 * partition join, leave, merge, split 등은
	 * partitionManager leader 에 의해 정해진다.
	 * follower는 partitionManagerFsm 에 의해 판단한다. */
	pm *partitionManager.PartitionManager

	/* partitionManager 에 의해 결정된 Event를
	 * partition 에 적용하기 위한 용도. */
	partEventCh *partitionManager.Event
}

func NewServer() (*Server, error) {
	s := &Server{}

	s.setupPartitionManager()

	s.setupXsmember()

	go s.eventLoop()

	s.p, _ = partition.NewPartition()

	return s, nil
}

func (s *Server) setupPartitionManager() {
	s.partEventCh = make(chan *partitionManager.Event, 256)
	s.pm, _ = partitionManager.NewPartitionManager(s.partEventCh)
}

func (s *Server) setupXsmember() {
	s.eventCh = make(chan *xsmember.Event, 256)
	conf := xsmember.NewConfig(s.eventCh)

	s.xsmember, _ = xsmember.NewXsmember(conf)
}

func (s *Server) eventLoop() {
	for {
		select {
		case e := <-s.eventCh:
			s.handleMemberEvent(e)
		case e := <-s.partEventCh:
			s.handlePartitionManageEvent(e)
		}
	}
}

/* xsmemeber로 부터 받은 event를
 * partitionManager 에 전달*/
func (s *Server) handleMemberEvent(e *xsmember.Event) {
	switch e.Type() {
	case xsmember.TypeJoin:
		s.pm.AddPeer(e.Member())
	case xsmember.TypeLeave:
		s.pm.LeavePeer(e.Member())
	case xsmember.TypeUpdate:
		s.pm.UpdatePeer(e.Member())
	}
}

/* partitionManager 로 부터 partiotn 할당, 삭제,
 * 또는 memge, split 명령을 받는다. */
func (s *Server) handlePartitionManageEvent(e *partitionManager.Event) {
	switch e.Type() {
	case partitionManager.TypeJoin:
	case xsmember.TypeLeave:
	case xsmember.TypeMerge:
	case xsmember.TypeSplit:
	}
}

func (s *Server) Start() {
}
