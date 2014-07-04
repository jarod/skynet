package main

import (
	"bytes"
	proto "code.google.com/p/goprotobuf/proto"
	"encoding/json"
	skynet "github.com/jarod/skynet/skynet"
	skc "github.com/jarod/skynet/skynet/client"
	skn "github.com/jarod/skynet/skynet/net"
	"io"
	"log"
	"net"
	"sync"
)

var (
	appIdMap map[string]*App
	connMap  map[*net.TCPConn]*App
	appInfos map[string]*skc.AppInfo
	mutex    sync.Mutex
)

func init() {
	appIdMap = make(map[string]*App)
	connMap = make(map[*net.TCPConn]*App)
	appInfos = make(map[string]*skc.AppInfo)
}

type App struct {
	id   string
	conn *net.TCPConn
}

func onAppConnected(conn *net.TCPConn) {
	ac := NewApp(conn)
	connMap[conn] = ac
	for {
		p, err := skn.ParsePacket(conn)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		ac.dispatchAppPacket(p)
	}
	onAppDisconnected(ac)
	conn.Close()
}

func onAppDisconnected(ac *App) {
	delete(connMap, ac.conn)
	id := ac.id
	if id != "" {
		delete(appIdMap, id)
		m := &skynet.Pstring{
			Value: proto.String(id),
		}
		p, _ := skn.NewMessagePacket(0x0001, m)
		matrixClient.Write(p)
	}
	log.Printf("app: disconnected id=%s,addr=%s", id, ac.conn.RemoteAddr())
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
		ac.sendToApp(p)
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

func (ac *App) sendToApp(p *skn.Packet) {
	msg := new(skynet.AppMsg)
	proto.Unmarshal(p.Body, msg)
	mutex.Lock()
	defer mutex.Unlock()
	if c, ok := appIdMap[*msg.AppId]; ok {
		//log.Printf("Msg local appId:%s,head=%d\n", *msg.AppId, *msg.Head)
		c.Write(p)
	} else {
		//log.Printf("Msg remote appId:%s,head=%d\n", *msg.AppId, *msg.Head)
		matrixClient.Write(p)
	}
}
