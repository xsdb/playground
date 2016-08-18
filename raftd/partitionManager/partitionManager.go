package partitionManager

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

type PartitionManagerFSM struct {
}

type PartitionManagerPeerStore struct {
}

type PartitionManager struct {
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

func NewPartitionManager(leader string, port int, dir string,
	eventCh chan *Event) (*PartitionManager, error) {

	bindAddr := fmt.Sprintf(":%v", port)
	log.Printf("[PM] bind %v...\n", bindAddr)
	advertise, err := net.ResolveTCPAddr("tcp", bindAddr)
	if err != nil {
		return nil, err
	}

	log.Printf("[PM] new transport")
	logger := log.New(os.Stdout, "[transport]", log.LstdFlags)
	trans, err := raft.NewTCPTransportWithLogger(bindAddr, advertise,
		maxPool, timeout, logger)
	if err != nil {
		return nil, err
	}

	conf := raft.DefaultConfig()
	if leader == "" {
		conf.EnableSingleNode = true
	}

	fsm := &PartitionManagerFSM{}

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

	peerStore := &PartitionManagerPeerStore{}

	log.Printf("[PM] new raft")
	raft, err := raft.NewRaft(conf, fsm, logs, logs, snaps, peerStore, trans)
	if err != nil {
		return nil, err
	}

	log.Printf("[PM] new partition manager")
	pm := &PartitionManager{raft: raft,
		eventCh: eventCh}

	return pm, nil
}

func (fsm *PartitionManagerFSM) Apply(l *raft.Log) interface{} {
	log.Println("[Apply to FSM] Index=", l.Index)
	log.Println("[Apply to FSM] Term=", l.Term)
	log.Println("[Apply to FSM] Type=", l.Type)
	return &PmLogCommitResult{}
}

func (fsm *PartitionManagerFSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (fsm *PartitionManagerFSM) Restore(reader io.ReadCloser) error {
	return nil
}

func (ps *PartitionManagerPeerStore) Peers() ([]string, error) {
	return []string{}, nil
}

func (ps *PartitionManagerPeerStore) SetPeers([]string) error {
	return nil
}

func (pm *PartitionManager) AddPeer(m *xsmember.Member) error {
	log.Printf("[PM] add peer %v\n", m)

	future := pm.raft.AddPeer(m.Addr.String())
	if err := future.Error(); err != nil && err != raft.ErrKnownPeer {
		log.Printf("[ERR] [PartitionManager] failed to add raft peer: %v", err)
		return err
	} else if err == nil {
		log.Printf("[INFO] [PartitionManager] added raft peer: %v", m)
	}
	return nil
}

func (pm *PartitionManager) LeavePeer(m *xsmember.Member) error {
	future := pm.raft.RemovePeer(m.Addr.String())
	if err := future.Error(); err != nil && err != raft.ErrKnownPeer {
		log.Printf("[ERR] [PartitionManager] failed to remove raft peer: %v", err)
		return err
	} else if err == nil {
		log.Printf("[INFO] [PartitionManager] removed raft peer: %v", m)
	}
	return nil
}

func (e *Event) Type() EventType {
	return e.eventType
}

func (e *Event) Member() *xsmember.Member {
	return e.member
}
