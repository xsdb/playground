package server

import (
	"github.com/xsdb/playground/raftd/partition"
)

type Server struct {
	p *partition.Partition
}

func NewServer() (*Server, error) {
	return &Server{}, nil
}
