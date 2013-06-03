package main

import (
	"flag"
	"fmt"
	smc "github.com/jarod/skynet/skynet/matrix/client"
	"log"
	"net"
	"os"
)

const (
	Version = "0.4-130531"
)

var version = flag.Bool("version", false, "show skynet-agent version")

var matrix = flag.String("matrix", "127.0.0.1:1860", "address of matrix server")

var (
	matrixClient *smc.MatrixClient
)

func bindAgentServer() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:1890")
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
			log.Println(err)
			continue
		}

		go onClientConnected(conn)
	}
	listener.Close()
}

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("skynet-agent - %s\n", Version)
		os.Exit(0)
	}

	mc, err := smc.Dial(*matrix)
	if err != nil {
		log.Println(err)
	}
	go readMatrix(mc)

	bindAgentServer()
}
