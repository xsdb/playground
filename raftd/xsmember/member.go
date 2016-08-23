package xsmember

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/hashicorp/memberlist"
)

type Xsmember struct {
	list       *memberlist.Memberlist
	broadcasts *memberlist.TransmitLimitedQueue
	conf       *Config
}

type Member struct {
	Name   string
	Addr   net.IP
	Port   uint16
	Tags   map[string]string
	Status MemberStatus

	ProtocolMin uint8
	ProtocolMax uint8
	ProtocolCur uint8
	DelegateMin uint8
	DelegateMax uint8
	DelegateCur uint8
}

type MemberStatus int

const (
	StatusNone MemberStatus = iota
	StatusAlive
	StatusLeaving
	StatusLeft
	StatusFailed
)

type Event struct {
	eventType EventType
	member    *Member
}

type EventType int

const (
	TypeNone EventType = iota
	TypeJoin
	TypeLeave
	TypeUpdate
)

type Config struct {
	memberlistConf *memberlist.Config
	eventCh        chan *Event
	msgCh          chan *Message
	Tags           map[string]string
}

func NewXsmember(conf *Config, name string, port int) (*Xsmember, error) {
	xsmember := &Xsmember{
		conf: conf}

	conf.memberlistConf = memberlist.DefaultLocalConfig()
	conf.memberlistConf.Name = name
	conf.memberlistConf.BindPort = port
	conf.memberlistConf.AdvertisePort = port
	conf.memberlistConf.Events = &eventDelegate{xsmember: xsmember}
	conf.memberlistConf.Delegate = &delegate{xsmember: xsmember}

	list, err := memberlist.Create(conf.memberlistConf)
	if err != nil {
		return nil, err
	}
	xsmember.list = list

	xsmember.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return xsmember.list.NumMembers()
		},
		RetransmitMult: conf.memberlistConf.RetransmitMult,
	}

	return xsmember, nil
}

func (xsmember *Xsmember) Join(addrs []string) (int, error) {
	n, err := xsmember.list.Join(addrs)

	// Ask for members of the cluster
	for _, member := range xsmember.list.Members() {
		log.Printf("[Xsmember] member: %s %s\n", member.Name, member.Addr)
	}

	return n, err
}

func (xsmember *Xsmember) Update(tags map[string]string) error {
	xsmember.conf.Tags = tags

	err := xsmember.list.UpdateNode(500 * time.Millisecond)
	if err != nil {
		log.Printf("[Xsmember] update fail %v", err)
	}

	return err
}

func (xsmember *Xsmember) Broadcast(t msgType, msg interface{}) {
	buf, _ := encodeMessage(t, msg)
	log.Printf("%v", buf)
	b := broadcast{msg: buf}
	xsmember.broadcasts.QueueBroadcast(b)
}

func NewConfig(eventCh chan *Event, msgCh chan *Message) *Config {
	return &Config{eventCh: eventCh, msgCh: msgCh}
}

func (xsmember *Xsmember) decodeTags(meta []byte) map[string]string {
	var tags map[string]string

	err := json.Unmarshal(meta, &tags)
	if err != nil {
		log.Printf("[Xsmember] decodeTags: %v", err)
	}

	return tags
}

func (xsmember *Xsmember) encodeTags(tags map[string]string) []byte {
	meta, err := json.Marshal(tags)
	if err != nil {
		log.Printf("[Xsmember] encodeTags: %v", err)
	}

	return meta
}

func (xsmember *Xsmember) handleMemberJoin(m *Member) {
	e := &Event{
		eventType: TypeJoin,
		member:    m}

	xsmember.conf.eventCh <- e
}

func (xsmember *Xsmember) handleMemberLeave(m *Member) {
	e := &Event{
		eventType: TypeLeave,
		member:    m}

	xsmember.conf.eventCh <- e
}

func (xsmember *Xsmember) handleMemberUpdate(m *Member) {
	e := &Event{
		eventType: TypeUpdate,
		member:    m}

	xsmember.conf.eventCh <- e
}

func (e *Event) Type() EventType {
	return e.eventType
}

func (e *Event) Member() *Member {
	return e.member
}
