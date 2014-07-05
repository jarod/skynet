package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var VERSION = "0.9-SNAPSHOT"

var version = flag.Bool("version", false, "show skynet-agent version")

var matrix = flag.String("matrix", "127.0.0.1:1860", "address of matrix server")

var optTcpAddr = flag.String("tcp", ":1890", "address to serve tcp")
var optHttpAddr = flag.String("http", ":1891", "address to serve http")

var (
	tcpServer    *TcpServer
	matrixClient *MatrixClient
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("skynet-agent - %s\n", VERSION)
		os.Exit(0)
	}
	var err error
	matrixClient, err = DialMatrix()
	if err != nil {
		log.Println("MatrixClient: ", err)
		return
	}
	httpServer := NewHttpServer()
	go httpServer.ListenAndServe(*optHttpAddr)
	tcpServer = NewTcpServer()
	tcpServer.ListenAndServe(*optTcpAddr)
}
