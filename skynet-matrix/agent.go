package main

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"github.com/jarod/skynet/skynet"
	skc "github.com/jarod/skynet/skynet/client"
	snet "github.com/jarod/skynet/skynet/net"
	"io"
	"log"
	"net"
	"sync"
)

var (
	agents []*Agent
	mutex  sync.Mutex
	apps   map[string]*skc.AppInfo // id->info
)

func init() {

}

func findAgentByIP(ip string) *Agent {
	for _, a := range agents {
		if a.remoteIp == ip {
			return a
		}
	}
	return nil
}

type Agent struct {
	conn     *net.TCPConn
	remoteIp string
}

func onAgentConnected(conn *net.TCPConn) {
	ag := NewAgent(conn)

	mutex.Lock()
	agents = append(agents, ag)
	mutex.Unlock()

	log.Printf("Agent connected %s\n", ag.RemoteIp())
	for {
		p, err := snet.ParsePacket(conn)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		ag.dispatchAgentPacket(p)
	}
	ag.onDisconnected()
	conn.Close()
}

func NewAgent(conn *net.TCPConn) (ac *Agent) {
	ac = new(Agent)
	ac.conn = conn
	return
}

func (a *Agent) onDisconnected() {
	mutex.Lock()
	for i, agent := range agents {
		if agent == a {
			agents[i], agents = agents[len(agents)-1], agents[:len(agents)-1]
			break
		}
	}
	mutex.Unlock()

	log.Printf("Agent disconnected. ip=%s\n", a.RemoteIp())
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

func (ac *Agent) Write(p *snet.Packet) {
	data := p.Encode()
	io.Copy(ac.conn, bytes.NewReader(data))
}

func (a *Agent) dispatchAgentPacket(p *snet.Packet) {
	switch p.Head {
	case 0x0000:
		a.updateAppInfo(p)
	case 0x0001:
		a.onAppDisconnected(p)
	default:
		a.broadcast(p)
	}
	log.Printf("dispatchAgentPacket - %v\n", p)
}

func (a *Agent) broadcast(p *snet.Packet) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, v := range agents {
		v.Write(p)
	}
}

func (a *Agent) updateAppInfo(p *snet.Packet) {
	info := new(skc.AppInfo)
	err := json.Unmarshal(p.Body, info)
	if err != nil {
		log.Println("updateAppInfo - ", err)
		return
	}
	info.Agent = a.conn.RemoteAddr().String()
	apps[info.Id] = info
	p.Body, err = json.Marshal(info)
	if err != nil {
		log.Println("updateAppInfo - ", err)
		return
	}
	a.broadcast(p)
}

func (a *Agent) onAppDisconnected(p *snet.Packet) {
	id := new(skynet.Pstring)
	err := proto.Unmarshal(p.Body, id)
	if err != nil {
		log.Println("onAppDisconnected - ", err)
		return
	}
	delete(apps, id.GetValue())
	a.broadcast(p)
}
