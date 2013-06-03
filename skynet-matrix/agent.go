package main

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"github.com/jarod/skynet/skynet"
	snet "github.com/jarod/skynet/skynet/net"
	"io"
	"log"
	"net"
	"sync"
)

var (
	ipMap       map[string]*Agent
	clientCount int
	mutex       *sync.Mutex
)

func init() {
	ipMap = make(map[string]*Agent)
	mutex = new(sync.Mutex)
}

type Agent struct {
	conn     *net.TCPConn
	clients  []int32
	remoteIp string
}

func onAgentConnected(conn *net.TCPConn) {
	ag := NewAgent(conn)

	mutex.Lock()
	ipMap[ag.RemoteIp()] = ag
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
	onAgentDisconnected(ag)
	conn.Close()
}

func onAgentDisconnected(a *Agent) {
	mutex.Lock()
	delete(ipMap, a.RemoteIp())
	mutex.Unlock()

	log.Printf("Agent disconnected. ip=%s\n", a.RemoteIp())
}

func NewAgent(conn *net.TCPConn) (ac *Agent) {
	ac = new(Agent)
	ac.conn = conn
	ac.clients = make([]int32, 0)
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

func (ac *Agent) Write(p *snet.Packet) {
	data := p.Encode()
	io.Copy(ac.conn, bytes.NewReader(data))
}

func (a *Agent) dispatchAgentPacket(p *snet.Packet) {
	switch p.Head {
	case 0x0000:
		a.registerClient(p)
	case 0x0001:
		a.onClientDisconnected(p)
	default:
		a.broadcast(p)
	}
	log.Printf("dispatchAgentPacket - %v\n", p)
}

func (a *Agent) broadcast(p *snet.Packet) {
	for _, v := range ipMap {
		v.Write(p)
	}
}

func (a *Agent) addClient(id int32) {
	a.clients = append(a.clients, id)
	clientCount++
}

func (a *Agent) delClient(id int32) {
	for i, v := range a.clients {
		if v == id {
			a.clients = append(a.clients[:i], a.clients[i+1:]...)
			clientCount--
		}
	}
}

func (a *Agent) registerClient(p *snet.Packet) {
	id := new(skynet.Psint32)
	proto.Unmarshal(p.Body, id)
	a.broadcast(p)
	a.addClient(id.GetValue())
}

func (a *Agent) onClientDisconnected(p *snet.Packet) {
	id := new(skynet.Psint32)
	proto.Unmarshal(p.Body, id)
	a.delClient(id.GetValue())

	a.broadcast(p)
}
