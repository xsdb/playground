package partitionMap

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/satori/go.uuid"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"

	"github.com/xsdb/playground/raftd/xsmember"
)

const (
	maxPool           = 128
	timeout           = time.Second * 3
	kfactor           = 2
	num_partition     = 8
	num_site_per_node = 1
)

type PmLogCommitResult struct {
	result int
}

type PartitionMapFSM struct {
	partMap *PartitionMap
	db      *leveldb.DB
}

type PartitionMapPeerStore struct {
}

type PartitionMap struct {
	raft    *raft.Raft
	fsm     *PartitionMapFSM
	eventCh chan *Event
}

type Event struct {
}

func NewPartitionMap(leader string, port int, dir string, eventCh chan *Event) (*PartitionMap, error) {

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

	db, err := leveldb.OpenFile(dir+"partMap/fsm", nil)
	if err != nil {
		return nil, err
	}

	fsm := &PartitionMapFSM{db: db}

	logs, err := raftboltdb.NewBoltStore(dir + "/partMap/logs")
	if err != nil {
		return nil, err
	}

	logger = log.New(os.Stdout, "[snapshot]", log.LstdFlags)
	snaps, err := raft.NewFileSnapshotStoreWithLogger(dir+"/pm/", 1, logger)
	if err != nil {
		return nil, err
	}

	peerStore := &PartitionMapPeerStore{}

	raft, err := raft.NewRaft(conf, fsm, logs, logs, snaps, peerStore, trans)
	if err != nil {
		return nil, err
	}

	pm := &PartitionMap{raft: raft,
		fsm:     fsm,
		eventCh: eventCh}

	fsm.partMap = pm

	return pm, nil
}

func (fsm *PartitionMapFSM) Apply(l *raft.Log) interface{} {
	log.Println("[Apply to FSM] Index=", l.Index)
	log.Println("[Apply to FSM] Term=", l.Term)
	log.Println("[Apply to FSM] Type=", l.Type)

	t := msgType(l.Data[0])
	switch t {
	case MsgTypeRegisterPartition:
		var partitionInfo PartitionInfo
		if err := decodeMessage(l.Data[1:], &partitionInfo); err != nil {
			log.Printf("[Apply to FSM] fail. %v", err)
			break
		}
		fsm.setPartitionInfo(&partitionInfo)
	case MsgTypeRegisterNodeInfo:
		var nodeInfo NodeInfo
		if err := decodeMessage(l.Data[1:], &nodeInfo); err != nil {
			log.Printf("[Apply to FSM] fail. %v", err)
			break
		}
		fsm.setNodeInfo(&nodeInfo)
	case MsgTypeRegisterSiteInfo:
		var siteInfo SiteInfo
		if err := decodeMessage(l.Data[1:], &siteInfo); err != nil {
			log.Printf("[Apply to FSM] fail. %v", err)
			break
		}
		fsm.setSiteInfo(&siteInfo)
	}

	select {
	case fsm.partMap.eventCh <- &Event{}:
	default:
	}

	return &PmLogCommitResult{}
}

func (fsm *PartitionMapFSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (fsm *PartitionMapFSM) Restore(reader io.ReadCloser) error {
	return nil
}

func (fsm *PartitionMapFSM) Get(key []byte) ([]byte, error) {
	return fsm.db.Get(key, nil)
}

func (fsm *PartitionMapFSM) Put(key []byte, val []byte) error {
	return fsm.db.Put(key, val, nil)
}

func (fsm *PartitionMapFSM) NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {
	return fsm.db.NewIterator(slice, ro)
}

func (ps *PartitionMapPeerStore) Peers() ([]string, error) {
	return []string{}, nil
}

func (ps *PartitionMapPeerStore) SetPeers([]string) error {
	return nil
}

func (pm *PartitionMap) AddPeer(name string, addr string, port int, dir string) error {
	future := pm.raft.AddPeer(fmt.Sprintf("%v:%v", addr, port))
	if err := future.Error(); err != nil && err != raft.ErrKnownPeer {
		log.Printf("[ERR] [PartitionMap] failed to add raft peer: %v", err)
		return err
	} else if err == nil {
		log.Printf("[INFO] [PartitionMap] added raft peer: %v:%v", addr, port)
		if pm.State() == raft.Leader {
			pm.RegisterNode(name, addr, port, dir)
		}
	}
	return nil
}

func (pm *PartitionMap) LeavePeer(addr string) error {
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

func (pm *PartitionMap) WaitState(state raft.RaftState) {
	for {
		if pm.State() == state {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (pm *PartitionMap) State() raft.RaftState {
	return pm.raft.State()
}

func (pm *PartitionMap) Get(key []byte) ([]byte, error) {
	return pm.fsm.Get(key)
}

func (pm *PartitionMap) NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {
	return pm.fsm.NewIterator(slice, ro)
}

func (pm *PartitionMap) InitPartitionTable() error {
	for i := 0; i < num_partition; i++ {
		err := pm.AppendPartition(i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pm *PartitionMap) RegisterNode(name string, addr string, port int, dir string) error {
	n, err := pm.AppendNodeInfo(name, fmt.Sprintf("%v:%v", addr, port))
	if err != nil {
		return err
	}

	log.Printf("NodeInfo[%v] Added", n.Nid)

	for i := 0; i < num_site_per_node; i++ {
		log.Printf("AppendSite")
		p, err := pm.AppendSiteInfo(n.Nid, port+1+i, dir)
		if err != nil {
			return err
		}

		log.Printf("SiteInfo[%v] Added", p.Pid)
	}

	log.Printf("AddPeer Complete")

	return nil
}

func (pm *PartitionMap) AppendPartition(pid int) error {
	m := &PartitionInfo{Pid: pid,
		Kfactor: kfactor,
	}

	if err := pm.appendMessage(MsgTypeRegisterPartition, m); err != nil {
		return err
	}

	return nil
}

func (pm *PartitionMap) AppendNodeInfo(name string, addr string) (*NodeInfo, error) {
	m := &NodeInfo{Nid: pm.getNewNid(),
		Name: name,
		Addr: addr,
	}

	if err := pm.appendMessage(MsgTypeRegisterNodeInfo, m); err != nil {
		return nil, err
	}

	return m, nil
}

func (pm *PartitionMap) AppendSiteInfo(nid uuid.UUID, port int, path string) (*SiteInfo, error) {
	pid := pm.getNewPid()
	p := &SiteInfo{Pid: pid,
		Nid:  nid,
		Port: port,
		Path: fmt.Sprintf("%v/store/%v", path, pid),
	}

	if err := pm.appendMessage(MsgTypeRegisterSiteInfo, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (pm *PartitionMap) getNewNid() uuid.UUID {
	return uuid.NewV4()
}

func (pm *PartitionMap) getNewPid() uuid.UUID {
	return uuid.NewV4()
}

func (pm *PartitionMap) appendMessage(t msgType, m interface{}) error {

	b, err := encodeMessage(t, m)
	if err != nil {
		return err
	}

	resp := pm.raft.Apply(b, 500*time.Millisecond)
	if err = resp.Error(); err != nil {
		return err
	}

	return nil
}

func (fsm *PartitionMapFSM) setPartitionInfo(p *PartitionInfo) error {
	for i := 0; i < p.Kfactor; i++ {
		key := fmt.Sprintf("/partition/%v/%v", p.Pid, i)
		val := fmt.Sprintf("none")

		log.Printf("set %v : %v", key, val)

		err := fsm.Put([]byte(key), []byte(val))
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (fsm *PartitionMapFSM) setNodeInfo(n *NodeInfo) error {
	key := fmt.Sprintf("/node/%v", n.Name)
	val := fmt.Sprintf("%v", n.Nid)

	log.Printf("set %v : %v", key, val)

	err := fsm.Put([]byte(key), []byte(val))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (fsm *PartitionMapFSM) setSiteInfo(p *SiteInfo) error {
	key := fmt.Sprintf("/site/%v/node", p.Pid)
	val := fmt.Sprintf("%v", p.Nid)

	log.Printf("set %v : %v", key, val)

	err := fsm.Put([]byte(key), []byte(val))
	if err != nil {
		log.Println(err)
		return err
	}
	key = fmt.Sprintf("/site/%v/path", p.Pid)
	val = fmt.Sprintf("%v", p.Path)

	log.Printf("set %v : %v", key, val)

	err = fsm.Put([]byte(key), []byte(val))
	if err != nil {
		log.Println(err)
		return err
	}

	key = fmt.Sprintf("/site/%v/port", p.Pid)
	val = fmt.Sprintf("%v", p.Port)

	log.Printf("set %v : %v", key, val)

	err = fsm.Put([]byte(key), []byte(val))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
