package main

import (
	"bytes"
	"encoding/json"
	skc "github.com/jarod/skynet/skynet/client"
	skn "github.com/jarod/skynet/skynet/net"
	"io"
	"log"
	"net"
)

var (
	appInfos map[string]*skc.AppInfo
)

func init() {
	appInfos = make(map[string]*skc.AppInfo)
}

type App struct {
	id   string
	conn *net.TCPConn
}

func NewApp(conn *net.TCPConn) (ac *App) {
	ac = new(App)
	ac.conn = conn
	return
}

func (ac *App) Write(p *skn.Packet) {
	data := p.Encode()
	io.Copy(ac.conn, bytes.NewReader(data))
}

func (ac *App) dispatchAppPacket(p *skn.Packet) {
	switch p.Head {
	case 0x0000:
		ac.updateAppInfo(p)
	case 0x0010:
		tcpServer.SendToApp(p)
	default:
		matrixClient.Write(p)
	}
	//log.Printf("dispatchAppPacket %v\n", p)
}

func (a *App) updateAppInfo(p *skn.Packet) {
	info := new(skc.AppInfo)
	err := json.Unmarshal(p.Body, info)
	if err != nil {
		log.Println("app: updateAppInfo - ", err)
		return
	}
	a.id = info.Id
	matrixClient.Write(p)
}
