package main

import (
	"flag"
	"log"
	"os"

	"github.com/xsdb/playground/raftd/server"
)

func main() {
	hostname, _ := os.Hostname()

	leaderPtr := flag.String("leader", "", "a string")
	portPtr := flag.Int("port", 50051, "an int")
	dirPtr := flag.String("dir", "./tmp/", "a string")
	namePtr := flag.String("name", hostname, "a string")

	flag.Parse()

	log.Printf("leader: %v", *leaderPtr)
	log.Printf("port:   %v", *portPtr)
	log.Printf("dir:    %v", *dirPtr)
	log.Printf("name:   %v", *namePtr)

	server, _ := server.NewServer(*leaderPtr, *portPtr, *dirPtr, *namePtr)
	server.Start()
}
