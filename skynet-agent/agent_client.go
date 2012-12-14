package main

import (
	"bytes"
	proto "code.google.com/p/goprotobuf/proto"
	skynet "github.com/jarod/skynet/skynet"
	snet "github.com/jarod/skynet/skynet/net"
	"io"
	"log"
	"net"
)

var (
	idMap   map[int32]*AgentClient
	connMap map[*net.TCPConn]*AgentClient
)

func init() {
	idMap = make(map[int32]*AgentClient)
	connMap = make(map[*net.TCPConn]*AgentClient)
}

type AgentClient struct {
	id   int32
	conn *net.TCPConn
}

func onClientConnected(conn *net.TCPConn) {
	ac := NewAgentClient(conn)
	connMap[conn] = ac
	for {
		p, err := snet.ParsePacket(conn)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		ac.dispatchClientPacket(p)
	}
	onClientDisconnected(ac)
	conn.Close()
}

func onClientDisconnected(ac *AgentClient) {
	delete(connMap, ac.conn)
	if ac.id >= 0 {
		delete(idMap, ac.id)
		m := &skynet.SInt32{
			Value: proto.Int32(ac.id),
		}
		p, _ := snet.NewMessagePacket(0x0001, m)
		matrixClient.Write(p)
	}
	log.Printf("Client disconnected id=%d,ip=%s", ac.id, ac.conn.RemoteAddr())
}

func NewAgentClient(conn *net.TCPConn) (ac *AgentClient) {
	ac = new(AgentClient)
	ac.id = -1
	ac.conn = conn
	return
}

func (ac *AgentClient) Write(p *snet.Packet) {
	data := p.Encode()
	io.Copy(ac.conn, bytes.NewReader(data))
}

func (ac *AgentClient) dispatchClientPacket(p *snet.Packet) {
	switch p.Head {
	case 0x0000:
		ac.registerClient(p)
	}
	log.Printf("dispatchClientPacket %v\n", p)
}

func (ac *AgentClient) registerClient(p *snet.Packet) {
	id := new(skynet.SInt32)
	proto.Unmarshal(p.Body, id)
	ac.id = id.GetValue()
	idMap[id.GetValue()] = ac
	matrixClient.Write(p)
	log.Printf("New client id=%d,ip=%s", ac.id, ac.conn.RemoteAddr())
}
