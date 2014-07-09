package main

import (
	"flag"
	"fmt"
	sklog "github.com/jarod/skynet/skynet/log"
	"os"
)

var VERSION = "0.9-SNAPSHOT"

var version = flag.Bool("version", false, "show skynet-matrix version")
var optLogFile = flag.String("log", "", "log file location, rotate on SIGUSR1")

var optTcpAddr = flag.String("tcp", ":1860", "address to serve tcp")
var optHttpAddr = flag.String("http", ":1861", "address to serve http")

var (
	httpServer *HttpServer
	tcpServer  *TcpServer
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("skynet-matrix - %s\n", VERSION)
		os.Exit(0)
	}

	sklog.RegisterRotate(*optLogFile)

	tcpServer = NewTcpServer()
	go tcpServer.ListenAndServe(*optTcpAddr)

	httpServer = NewHttpServer()
	httpServer.ListenAndServe(*optHttpAddr)
}
