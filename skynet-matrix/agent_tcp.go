package main

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"github.com/jarod/skynet/skynet"
	skn "github.com/jarod/skynet/skynet/net"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type Agent struct {
	conn     *net.TCPConn
	remoteIp string
}

func NewAgent(conn *net.TCPConn) (ac *Agent) {
	ac = new(Agent)
	ac.conn = conn
	return
}

func (a *Agent) RemoteIp() string {
	if len(a.remoteIp) < 1 {
		tcpAddr, ok := a.conn.RemoteAddr().(*net.TCPAddr)
		if !ok {
			return ""
		}
		a.remoteIp = tcpAddr.IP.String()
	}
	return a.remoteIp
}

func (ac *Agent) Write(p *skn.Packet) {
	data := p.Encode()
	io.Copy(ac.conn, bytes.NewReader(data))
}

func (a *Agent) RequestAppInfos() {
	if !atomic.CompareAndSwapUint32(&appInfoLoaded, 0, 1) {
		return
	}
	p := skn.NewPacket(uint16(skynet.SkynetMsg_SM_AGENT_FIND_APPS), []byte("^.*$"))
	a.Write(p)
}

func (a *Agent) dispatchAgentPacket(p *skn.Packet) {
	switch skynet.SkynetMsg(p.Head) {
	case skynet.SkynetMsg_SM_APP_INFO:
		a.updateAppInfo(p)
	case skynet.SkynetMsg_SM_APP_DISCONNECTED:
		a.onAppDisconnected(p)
	case skynet.SkynetMsg_SM_AGENT_FIND_APPS:
		a.responseAppInfos(p)
	default:
		tcpServer.Broadcast(p)
	}
	log.Printf("dispatchAgentPacket - %v\n", p)
}

func (a *Agent) updateAppInfo(p *skn.Packet) {
	info := new(skynet.AppInfo)
	err := json.Unmarshal(p.Body, info)
	if err != nil {
		log.Println("updateAppInfo - ", err)
		return
	}
	info.Agent = a.conn.RemoteAddr().String()
	appInfos[info.Id] = info
	p.Body, err = json.Marshal(info)
	if err != nil {
		log.Println("updateAppInfo - ", err)
		return
	}
	tcpServer.Broadcast(p)
}

func (a *Agent) onAppDisconnected(p *skn.Packet) {
	id := new(skynet.Pstring)
	err := proto.Unmarshal(p.Body, id)
	if err != nil {
		log.Println("onAppDisconnected - ", err)
		return
	}
	delete(appInfos, id.GetValue())
	tcpServer.Broadcast(p)
}

func (a *Agent) responseAppInfos(p *skn.Packet) {
	var infos []*skynet.AppInfo
	err := json.Unmarshal(p.Body, &infos)
	if err != nil {
		log.Println("responseAppInfos - ", err)
		return
	}
	log.Println("Load app info from agent", infos)
	for _, info := range infos {
		_, exists := appInfos[info.Id]
		if exists {
			continue
		}
		appInfos[info.Id] = info
	}
}
