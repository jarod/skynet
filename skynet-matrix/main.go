package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var VERSION = "0.6-SNAPSHOT"

var version = flag.Bool("version", false, "show skynet-matrix version")
var optTcpAddr = flag.String("tcp", ":1860", "address to serve tcp")
var optHttpAddr = flag.String("http", ":1880", "address to serve http")

var (
	httpServer *MatrixHttpServer
)

func bindMatrixServer() {
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
		fmt.Printf("skynet-matrix - %s\n", VERSION)
		os.Exit(0)
	}
	go bindMatrixServer()

	httpServer = NewMatrixHttpServer()
	httpServer.Startup(*optHttpAddr)
}
