package server

import (
	"fmt"
	"log"
	"net"

	"github.com/xsdb/playground/raftd/partition"
	"github.com/xsdb/playground/raftd/partitionManager"
	"github.com/xsdb/playground/raftd/xsmember"
	"github.com/xsdb/playground/raftd/xsproto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	port int
	dir  string
	name string

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
	partEventCh chan *partitionManager.Event
}

func NewServer(leader string, port int, dir string, name string) (*Server, error) {
	s := &Server{port: port, dir: dir, name: name}

	log.Printf("setup PartitionManager\n")
	err := s.setupPartitionManager(leader, port+1)
	if err != nil {
		return nil, err
	}

	log.Printf("setup xsmember\n")
	err = s.setupXsmember(leader, port+2)
	if err != nil {
		return nil, err
	}

	log.Printf("run eventLoop\n")
	go s.eventLoop()

	/*
		log.Printf("create partition\n")
		s.p, err = partition.NewPartition(port + 2)
		if err != nil {
			return nil, err
		}
	*/

	return s, nil
}

func (s *Server) setupPartitionManager(leader string, port int) error {
	var err error

	s.partEventCh = make(chan *partitionManager.Event, 256)
	s.pm, err = partitionManager.NewPartitionManager(leader, port,
		s.dir, s.partEventCh)
	if err != nil {
		log.Printf("setup partition manager fail %v", err)
		return err
	}

	return nil
}

func (s *Server) setupXsmember(leader string, port int) error {
	var err error

	s.eventCh = make(chan *xsmember.Event, 256)
	conf := xsmember.NewConfig(s.eventCh)

	s.xsmember, err = xsmember.NewXsmember(conf, s.name, port)
	if err != nil {
		log.Printf("setup xsmember fail %v", err)
		return err
	}

	if leader != "" {
		_, err := s.xsmember.Join([]string{leader})
		if err != nil {
			log.Printf("Failed to join cluster: %v", err)
			return err
		}
	}

	return nil
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
	log.Printf("recv member Event %v, %v\n", e.Type(), e.Member())
	switch e.Type() {
	case xsmember.TypeJoin:
		s.pm.AddPeer(e.Member())
	case xsmember.TypeLeave:
		s.pm.LeavePeer(e.Member())
	case xsmember.TypeUpdate:
		//s.pm.UpdatePeer(e.Member())
	}
}

/* partitionManager 로 부터 partiotn 할당, 삭제,
 * 또는 memge, split 명령을 받는다. */
func (s *Server) handlePartitionManageEvent(e *partitionManager.Event) {
	switch e.Type() {
	case partitionManager.TypeJoin:
	}
}

func (s *Server) Get(ctx context.Context, in *xsproto.Request) (*xsproto.Response, error) {
	return &xsproto.Response{ReturnCode: 0, Value: "ok"}, nil
}

func (s *Server) Put(ctx context.Context, in *xsproto.Request) (*xsproto.Response, error) {
	return &xsproto.Response{ReturnCode: 0, Value: "ok"}, nil
}

func (s *Server) Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcserver := grpc.NewServer()
	xsproto.RegisterXsProtoServer(grpcserver, s)
	grpcserver.Serve(lis)
}
