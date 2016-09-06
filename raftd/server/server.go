package server

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/raft"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/xsdb/playground/raftd/partition"
	"github.com/xsdb/playground/raftd/partitionMap"
	"github.com/xsdb/playground/raftd/xsmember"
	"github.com/xsdb/playground/raftd/xsproto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Config struct {
	Clientport  int
	Memberport  int
	Partmapport int
	Dir         string
	Name        string
	Leader      string
}

type Server struct {
	conf *Config

	/* for cluster member manage */
	xsmember *xsmember.Xsmember

	/* for recv cluster member event */
	eventCh chan *xsmember.Event

	/* for recv cluster Msg*/
	msgCh chan *xsmember.Message

	/* for partition manage
	 * partition join, leave, merge, split 등은
	 * partitionMap leader 에 의해 정해진다.
	 * follower는 partitionMapFsm 에 의해 판단한다. */
	pm *partitionMap.PartitionMap

	/* partitionMap 에 의해 결정된 Event를
	 * partition 에 적용하기 위한 용도. */
	partEventCh chan *partitionMap.Event

	/* local partition map */
	localPartitionMap map[string]*partition.Partition
}

func NewServer(conf *Config) (*Server, error) {
	s := &Server{conf: conf}

	s.localPartitionMap = make(map[string]*partition.Partition)
	log.Printf("setup PartitionManager\n")
	err := s.setupPartitionManager()
	if err != nil {
		return nil, err
	}

	err = s.setupPartition()
	if err != nil {
		return nil, err
	}

	log.Printf("setup xsmember\n")
	err = s.setupXsmember()
	if err != nil {
		return nil, err
	}

	go s.eventLoop()

	if s.conf.Leader != "" {
		/* Join Memberlist */
		_, err := s.xsmember.Join([]string{s.conf.Leader})
		if err != nil {
			log.Printf("Failed to join cluster: %v", err)
		} else {
			/* Broadcast PartitionMap AddPeer Req */
			m := &xsmember.MsgAddPeer{Name: conf.Name,
				Addr: "localhost",
				Dir:  conf.Dir,
				Port: s.conf.Partmapport}
			s.xsmember.Broadcast(xsmember.MsgTypeAddPeer, m)
		}

		/* wait follower state */
		s.pm.WaitState(raft.Follower)
	} else {
		/* wait leader state */
		s.pm.WaitState(raft.Leader)

		s.pm.InitPartitionTable()
		s.pm.RegisterNode(conf.Name, "localhost", s.conf.Partmapport, s.conf.Dir)
	}

	return s, nil
}

func (s *Server) setupPartitionManager() error {
	var err error

	s.partEventCh = make(chan *partitionMap.Event, 256)
	s.pm, err = partitionMap.NewPartitionMap(s.conf.Leader,
		s.conf.Partmapport,
		s.conf.Dir,
		s.partEventCh)
	if err != nil {
		log.Printf("setup partition manager fail %v", err)
		return err
	}

	return nil
}

func (s *Server) setupXsmember() error {
	var err error

	s.eventCh = make(chan *xsmember.Event, 256)
	s.msgCh = make(chan *xsmember.Message, 256)
	conf := xsmember.NewConfig(s.eventCh, s.msgCh)

	s.xsmember, err = xsmember.NewXsmember(conf, s.conf.Name, s.conf.Memberport)
	if err != nil {
		log.Printf("setup xsmember fail %v", err)
		return err
	}

	return nil
}

func (s *Server) setupPartition() error {
	nid, err := s.pm.Get([]byte(fmt.Sprintf("/node/%v", s.conf.Name)))
	if err != nil {
		return nil
	}
	log.Printf("nodeID:%v", string(nid))

	iter := s.pm.NewIterator(util.BytesPrefix([]byte("/partition/")), nil)
	for iter.Next() {
		key := iter.Key()
		columns := strings.Split(string(key), "/")

		if len(columns) != 4 {
			continue
		}

		if columns[3] != "node" {
			continue
		}

		value := iter.Value()
		if string(nid) != string(value) {
			continue
		}

		pid := string(columns[2])
		if _, ok := s.localPartitionMap[pid]; ok == true {
			log.Printf("pid[%v] exist", pid)
			continue
		}

		path, err := s.pm.Get([]byte(fmt.Sprintf("/partition/%v/path", pid)))
		if err != nil {
			break
		}
		log.Printf("path : %v", string(path))

		port, err := s.pm.Get([]byte(fmt.Sprintf("/partition/%v/port", pid)))
		if err != nil {
			break
		}
		log.Printf("port : %v", string(port))

		p, err := partition.NewPartition(string(path), string(port))
		if err != nil {
			log.Printf("Partition create fail %v", err)
			break
		}
		log.Printf("NewPartition %v %v", path, port)

		s.localPartitionMap[pid] = p
	}
	iter.Release()
	err = iter.Error()

	return err
}

func (s *Server) eventLoop() {
	for {
		select {
		case e := <-s.eventCh:
			s.handleMemberEvent(e)
		case m := <-s.msgCh:
			s.handleMessage(m)
		case e := <-s.partEventCh:
			s.handlePartitionMapEvent(e)
		}
	}
}

func (s *Server) handleMemberEvent(e *xsmember.Event) {
	log.Printf("recv member Event %v, %v\n", e.Type(), e.Member())
	switch e.Type() {
	case xsmember.TypeJoin:
	case xsmember.TypeLeave:
	case xsmember.TypeUpdate:
	}
}

func (s *Server) handleMessage(m *xsmember.Message) {
	switch m.Type {
	case xsmember.MsgTypeAddPeer:
		msg := m.Body.(*xsmember.MsgAddPeer)
		s.pm.AddPeer(msg.Name, msg.Addr, msg.Port, msg.Dir)
	}
}

func (s *Server) handlePartitionMapEvent(e *partitionMap.Event) {
	s.setupPartition()
}

func (s *Server) Get(ctx context.Context, in *xsproto.Request) (*xsproto.Response, error) {
	return &xsproto.Response{ReturnCode: 0, Value: "ok"}, nil
}

func (s *Server) Put(ctx context.Context, in *xsproto.Request) (*xsproto.Response, error) {
	return &xsproto.Response{ReturnCode: 0, Value: "ok"}, nil
}

func (s *Server) Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.conf.Clientport))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcserver := grpc.NewServer()
	xsproto.RegisterXsProtoServer(grpcserver, s)
	grpcserver.Serve(lis)
}
