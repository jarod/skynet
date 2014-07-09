package main

import (
	"flag"
	"fmt"
	sklog "github.com/jarod/skynet/skynet/log"
	"log"
	"os"
)

var VERSION = "0.9-SNAPSHOT"

var version = flag.Bool("version", false, "show skynet-agent version")
var optLogFile = flag.String("log", "", "log file location, rotate on SIGUSR1")

var optMatrixAddr = flag.String("matrix", "127.0.0.1:1860", "address of matrix server")
var optMatrixUrl = flag.String("matrix-url", "http://127.0.0.1:1861/", "url of matrix http server")

var optTcpAddr = flag.String("tcp", ":1890", "address to serve tcp")
var optHttpAddr = flag.String("http", ":1891", "address to serve http")

var (
	tcpServer        *TcpServer
	matrixClient     *MatrixClient
	matrixHttpClient *MatrixHttpClient
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("skynet-agent - %s\n", VERSION)
		os.Exit(0)
	}

	sklog.RegisterRotate(*optLogFile)

	var err error
	matrixClient, err = DialMatrix(*optMatrixAddr)
	if err != nil {
		log.Println("MatrixClient: ", err)
		return
	}
	matrixHttpClient = NewMatrixHttpClient(*optMatrixUrl)
	fetchAppInfos()

	httpServer := NewHttpServer()
	go httpServer.ListenAndServe(*optHttpAddr)
	tcpServer = NewTcpServer()
	tcpServer.ListenAndServe(*optTcpAddr)
}

func fetchAppInfos() {
	infos := matrixHttpClient.FindApps(".*")
	for _, info := range infos {
		appInfos[info.Id] = info
	}
	log.Println("Fetch initial app info from matrix:", appInfos)
}
