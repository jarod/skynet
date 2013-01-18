package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	Version = "0.3-130117"
)

var version = flag.Bool("version", false, "show skynet-matrix version")

func bindMatrixServer() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:1860")
	if err != nil {
		log.Fatalln(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Listening on %s", addr.String())

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("AcceptTCP: %s\n", err)
			continue
		}
		go onAgentConnected(conn)
	}
	listener.Close()
}

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("skynet-matrix - %s\n", Version)
		os.Exit(0)
	}
	bindMatrixServer()
}
