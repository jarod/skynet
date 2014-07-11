package main

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/jarod/skynet/skynet"
	skn "github.com/jarod/skynet/skynet/net"
	"io"
	"log"
	"net"
	"sync"
)

type TcpServer struct {
	sync.Mutex
	appIdMap map[string]*App
	connMap  map[*net.TCPConn]*App
}

func NewTcpServer() *TcpServer {
	return &TcpServer{
		appIdMap: make(map[string]*App),
		connMap:  make(map[*net.TCPConn]*App),
	}
}

func (as *TcpServer) ListenAndServe(laddr string) {
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
			log.Println(err)
			continue
		}

		go as.onAppConnected(conn)
	}
	listener.Close()
}

func (as *TcpServer) onAppConnected(conn *net.TCPConn) {
	ac := NewApp(conn)
	as.connMap[conn] = ac
	log.Printf("New app connection %s\n", ac.conn.RemoteAddr())
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
	as.onAppDisconnected(ac)
	conn.Close()
}

func (as *TcpServer) onAppDisconnected(ac *App) {
	delete(as.connMap, ac.conn)
	id := ac.id
	if id != "" {
		delete(as.appIdMap, id)
		m := &skynet.Pstring{
			Value: proto.String(id),
		}
		p, _ := skn.NewMessagePacket(0x0001, m)
		matrixClient.Write(p)
	}
	log.Printf("app: disconnected id=%s,addr=%s", id, ac.conn.RemoteAddr())
}

func (as *TcpServer) BroadcastApps(p *skn.Packet) {
	for _, v := range as.connMap {
		v.Write(p)
	}
}

func (as *TcpServer) SendToApp(p *skn.Packet) {
	msg := new(skynet.AppMsg)
	err := proto.Unmarshal(p.Body, msg)
	if err != nil {
		log.Println("SendToApp - ", err)
		return
	}
	as.Lock()
	defer as.Unlock()
	if c, ok := as.appIdMap[*msg.AppId]; ok {
		//log.Printf("Msg local appId:%s,head=%d\n", *msg.AppId, *msg.Head)
		c.Write(p)
	} else {
		//log.Printf("Msg remote appId:%s,head=%d\n", *msg.AppId, *msg.Head)
		matrixClient.Write(p)
	}
}
