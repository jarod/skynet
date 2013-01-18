package main

import (
	"bytes"
	snet "github.com/jarod/skynet/skynet/net"
	"io"
	"log"
	"net"
)

var (
	connMap map[*net.TCPConn]*Agent
)

func init() {
	connMap = make(map[*net.TCPConn]*Agent)
}

type Agent struct {
	conn *net.TCPConn
}

func onAgentConnected(conn *net.TCPConn) {
	ag := NewAgent(conn)
	connMap[conn] = ag
	log.Printf("Agent connected %s\n", conn.RemoteAddr())
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

func onAgentDisconnected(ac *Agent) {
	delete(connMap, ac.conn)
	log.Printf("Agent disconnected ip=%s\n", ac.conn.RemoteAddr())
}

func NewAgent(conn *net.TCPConn) (ac *Agent) {
	ac = new(Agent)
	ac.conn = conn
	return
}

func (ac *Agent) Write(p *snet.Packet) {
	data := p.Encode()
	io.Copy(ac.conn, bytes.NewReader(data))
}

func (ac *Agent) dispatchAgentPacket(p *snet.Packet) {
	switch p.Head {
	case 0x0010:
		for _, v := range connMap {
			if ac != v {
				v.Write(p)
			}
		}
	default:
		for _, v := range connMap {
			v.Write(p)
		}
	}
	log.Printf("dispatchAgentPacket %v\n", p)
}
