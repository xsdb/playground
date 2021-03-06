package partition

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

const (
	maxPool = 128
	timeout = time.Second * 3
)

type LogCommitResult struct {
	result int
}

type PartitionFSM struct {
}

type PartitionPeerStore struct {
}

type Partition struct {
	raft *raft.Raft
}

func NewPartition(path string, port string) (*Partition, error) {
	os.MkdirAll(path, 0766)

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

	fsm := &PartitionFSM{}

	logs, err := raftboltdb.NewBoltStore(path + "/boltdb")
	if err != nil {
		return nil, err
	}

	logger = log.New(os.Stdout, "[snapshot]", log.LstdFlags)
	snaps, err := raft.NewFileSnapshotStoreWithLogger(path, 1, logger)
	if err != nil {
		return nil, err
	}

	peerStore := &PartitionPeerStore{}

	raft, err := raft.NewRaft(conf, fsm, logs, logs, snaps, peerStore, trans)
	if err != nil {
		return nil, err
	}

	partition := &Partition{raft: raft}

	return partition, nil
}

func (fsm *PartitionFSM) Apply(l *raft.Log) interface{} {
	log.Println("[Apply to FSM] Index=", l.Index)
	log.Println("[Apply to FSM] Term=", l.Term)
	log.Println("[Apply to FSM] Type=", l.Type)
	return &LogCommitResult{}
}

func (fsm *PartitionFSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (fsm *PartitionFSM) Restore(reader io.ReadCloser) error {
	return nil
}

func (ps *PartitionPeerStore) Peers() ([]string, error) {
	return []string{}, nil
}

func (ps *PartitionPeerStore) SetPeers([]string) error {
	return nil
}
