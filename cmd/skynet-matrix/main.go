package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/jarod/skynet/skynet"
	sklog "github.com/jarod/skynet/skynet/log"
)

var VERSION = "0.11-SNAPSHOT"

var version = flag.Bool("version", false, "show skynet-matrix version")
var optLogFile = flag.String("log", "", "log file location, reopen on signal SIGUSR1")

var optTcpAddr = flag.String("tcp", ":1860", "address to serve tcp")
var optHttpAddr = flag.String("http", ":1861", "address to serve http")

var (
	httpServer *HttpServer
	tcpServer  *TcpServer

	mutex    sync.Mutex
	appInfos map[string]*skynet.AppInfo // id->info

	appInfoLoaded uint32
)

func init() {
	appInfos = make(map[string]*skynet.AppInfo)
}

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("skynet-matrix - %s\n", VERSION)
		os.Exit(0)
	}

	sklog.SetLogFile(*optLogFile)

	tcpServer = NewTcpServer()
	go tcpServer.ListenAndServe(*optTcpAddr)

	httpServer = NewHttpServer()
	httpServer.ListenAndServe(*optHttpAddr)
}
