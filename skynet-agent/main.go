package main

import (
	"flag"
	"fmt"
	skmc "github.com/jarod/skynet/skynet/matrix/client"
	"log"
	"net"
	"os"
)

var VERSION = "0.9-SNAPSHOT"

var version = flag.Bool("version", false, "show skynet-agent version")

var matrix = flag.String("matrix", "127.0.0.1:1860", "address of matrix server")

var optTcpAddr = flag.String("tcp", ":1890", "address to serve tcp")

var (
	matrixClient *skmc.MatrixClient
)

func bindAgentServer() {
	addr, err := net.ResolveTCPAddr("tcp", *optTcpAddr)
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

		go onAppConnected(conn)
	}
	listener.Close()
}

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("skynet-agent - %s\n", VERSION)
		os.Exit(0)
	}

	mc, err := skmc.Dial(*matrix)
	if err != nil {
		log.Fatalln(err)
	}
	go readMatrix(mc)

	bindAgentServer()
}
