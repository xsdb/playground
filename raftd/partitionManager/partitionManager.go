package partitionManager

import (
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

const (
	bindAddr = "0.0.0.0"
	maxPool  = 128
	timeout  = time.Second * 3
)

type PmLogCommitResult struct {
	result int
}

type PartitionManagerFSM struct {
}

type PartitionManagerPeerStore struct {
}

type PartitionManager struct {
	raft *raft.Raft
}

func NewPartitionManager() (*PartitionManager, error) {
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

	fsm := &PartitionManagerFSM{}

	logs, err := raftboltdb.NewBoltStore("/tmp/pm/boltdb")
	if err != nil {
		return nil, err
	}

	logger = log.New(os.Stdout, "[snapshot]", log.LstdFlags)
	snaps, err := raft.NewFileSnapshotStoreWithLogger("/tmp/pm/", 1, logger)
	if err != nil {
		return nil, err
	}

	peerStore := &PartitionManagerPeerStore{}

	raft, err := raft.NewRaft(conf, fsm, logs, logs, snaps, peerStore, trans)
	if err != nil {
		return nil, err
	}

	pm := &PartitionManager{raft: raft}

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
