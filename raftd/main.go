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
	cportPtr := flag.Int("clientport", 50051, "an int")
	mportPtr := flag.Int("memberport", 50052, "an int")
	pportPtr := flag.Int("partmapport", 50053, "an int")
	dirPtr := flag.String("dir", "./tmp/", "a string")
	namePtr := flag.String("name", hostname, "a string")

	flag.Parse()

	log.Printf("leader: 	  %v", *leaderPtr)
	log.Printf("clientport:   %v", *cportPtr)
	log.Printf("memberport:   %v", *mportPtr)
	log.Printf("partmapport:  %v", *pportPtr)
	log.Printf("dir:    	  %v", *dirPtr)
	log.Printf("name:   	  %v", *namePtr)

	conf := &server.Config{Clientport: *cportPtr,
		Memberport:  *mportPtr,
		Partmapport: *pportPtr,
		Dir:         *dirPtr,
		Name:        *namePtr,
		Leader:      *leaderPtr}

	server, _ := server.NewServer(conf)
	server.Start()
}
