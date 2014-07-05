package main

import (
	skn "github.com/jarod/skynet/skynet/net"
	"io"
	"log"
	"net"
	"sync"
)

type TcpServer struct {
	sync.Mutex
	agents []*Agent
}

func NewTcpServer() *TcpServer {
	return &TcpServer{}
}

func (t *TcpServer) ListenAndServe(laddr string) {
	addr, err := net.ResolveTCPAddr("tcp", laddr)
	if err != nil {
		log.Fatalln(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Listening on %s", addr.String())

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("AcceptTCP: %s\n", err)
			continue
		}
		go t.onConnected(conn)
	}
	listener.Close()
}

func (t *TcpServer) onConnected(conn *net.TCPConn) {
	ag := NewAgent(conn)

	mutex.Lock()
	t.agents = append(t.agents, ag)
	mutex.Unlock()

	log.Printf("Agent connected %s\n", ag.RemoteIp())
	for {
		p, err := skn.ParsePacket(conn)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		ag.dispatchAgentPacket(p)
	}
	t.onDisconnected(ag)
	conn.Close()
}

func (t *TcpServer) onDisconnected(a *Agent) {
	t.Lock()
	ags := t.agents
	for i, agent := range ags {
		if agent == a {
			ags[i], ags = ags[len(ags)-1], ags[:len(ags)-1]
			break
		}
	}
	t.Unlock()

	log.Printf("Agent disconnected. ip=%s\n", a.RemoteIp())
}

func (t *TcpServer) FindAgentByIP(ip string) *Agent {
	t.Lock()
	defer t.Unlock()
	for _, a := range t.agents {
		if a.remoteIp == ip {
			return a
		}
	}
	return nil
}

func (t *TcpServer) Broadcast(p *skn.Packet) {
	t.Lock()
	defer t.Unlock()
	for _, v := range t.agents {
		v.Write(p)
	}
}
