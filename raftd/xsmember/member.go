package xsmember

import (
	"log"
	"net"

	"github.com/hashicorp/memberlist"
)

type Xsmember struct {
	list *memberlist.Memberlist
	conf *Config
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
}

func NewXsmember(conf *Config, name string, port int) (*Xsmember, error) {
	xsmember := &Xsmember{
		conf: conf}

	conf.memberlistConf = memberlist.DefaultLocalConfig()
	conf.memberlistConf.Name = name
	conf.memberlistConf.BindPort = port
	conf.memberlistConf.AdvertisePort = port
	conf.memberlistConf.Events = &eventDelegate{xsmember: xsmember}

	list, err := memberlist.Create(conf.memberlistConf)
	if err != nil {
		return nil, err
	}

	xsmember.list = list

	return xsmember, nil
}

func (xsmember *Xsmember) Join(addrs []string) (int, error) {
	n, err := xsmember.list.Join(addrs)

	// Ask for members of the cluster
	for _, member := range xsmember.list.Members() {
		log.Printf("Member: %s %s\n", member.Name, member.Addr)
	}

	return n, err
}

func NewConfig(c chan *Event) *Config {
	return &Config{eventCh: c}
}

func (xsmember *Xsmember) decodeTags(meta []byte) map[string]string {
	m := make(map[string]string)

	return m
}

func (xsmember *Xsmember) encodeTags(tags map[string]string) []byte {
	b := make([]byte, 8)

	return b
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
