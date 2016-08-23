package partitionMap

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"github.com/xsdb/playground/raftd/xsmember"
)

const (
	maxPool = 128
	timeout = time.Second * 3
)

type PmLogCommitResult struct {
	result int
}

type PartitionMapFSM struct {
}

type PartitionMapPeerStore struct {
}

type PartitionMap struct {
	raft    *raft.Raft
	eventCh chan *Event
}

type Event struct {
	eventType EventType
	member    *xsmember.Member
}

type EventType int

const (
	TypeNone EventType = iota
	TypeJoin
)

func NewPartitionMap(leader string, port int, dir string, eventCh chan *Event) (*PartitionMap, error) {

	log.Printf("[PM] new transport: %v", port)
	bindAddr := fmt.Sprintf(":%v", port)
	advertise, err := net.ResolveTCPAddr("tcp", bindAddr)
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "[transport]", log.LstdFlags)
	trans, err := raft.NewTCPTransportWithLogger(bindAddr, advertise, maxPool, timeout, logger)
	if err != nil {
		return nil, err
	}

	conf := raft.DefaultConfig()
	if leader == "" {
		conf.EnableSingleNode = true
	}

	fsm := &PartitionMapFSM{}

	log.Printf("[PM] new store")
	logs, err := raftboltdb.NewBoltStore(dir + "/pm/boltdb")
	if err != nil {
		return nil, err
	}

	log.Printf("[PM] new snapshot")
	logger = log.New(os.Stdout, "[snapshot]", log.LstdFlags)
	snaps, err := raft.NewFileSnapshotStoreWithLogger(dir+"/pm/", 1, logger)
	if err != nil {
		return nil, err
	}

	peerStore := &PartitionMapPeerStore{}

	log.Printf("[PM] new raft")
	raft, err := raft.NewRaft(conf, fsm, logs, logs, snaps, peerStore, trans)
	if err != nil {
		return nil, err
	}

	log.Printf("[PM] new partition manager")
	pm := &PartitionMap{raft: raft,
		eventCh: eventCh}

	return pm, nil
}

func (fsm *PartitionMapFSM) Apply(l *raft.Log) interface{} {
	log.Println("[Apply to FSM] Index=", l.Index)
	log.Println("[Apply to FSM] Term=", l.Term)
	log.Println("[Apply to FSM] Type=", l.Type)
	return &PmLogCommitResult{}
}

func (fsm *PartitionMapFSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (fsm *PartitionMapFSM) Restore(reader io.ReadCloser) error {
	return nil
}

func (ps *PartitionMapPeerStore) Peers() ([]string, error) {
	return []string{}, nil
}

func (ps *PartitionMapPeerStore) SetPeers([]string) error {
	return nil
}

func (pm *PartitionMap) AddPeer(addr string) error {
	log.Printf("[PM] add peer %v\n", addr)

	future := pm.raft.AddPeer(addr)
	if err := future.Error(); err != nil && err != raft.ErrKnownPeer {
		log.Printf("[ERR] [PartitionMap] failed to add raft peer: %v", err)
		return err
	} else if err == nil {
		log.Printf("[INFO] [PartitionMap] added raft peer: %v", addr)
	}
	return nil
}

func (pm *PartitionMap) LeavePeer(addr string) error {
	log.Printf("[PM] leave peer %v\n", addr)

	future := pm.raft.RemovePeer(addr)
	if err := future.Error(); err != nil && err != raft.ErrKnownPeer {
		log.Printf("[ERR] [PartitionMap] failed to remove raft peer: %v", err)
		return err
	} else if err == nil {
		log.Printf("[INFO] [PartitionMap] removed raft peer: %v", addr)
	}
	return nil
}

func (pm *PartitionMap) UpdatePeer(m *xsmember.Member) error {
	var err error

	return err
}

func (e *Event) Type() EventType {
	return e.eventType
}

func (e *Event) Member() *xsmember.Member {
	return e.member
}
