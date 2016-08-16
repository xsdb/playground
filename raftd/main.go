package main

import (
	"github.com/xsdb/playground/raftd/server"
)

func main() {
	server, _ := server.NewServer()
	server.Start()
}
