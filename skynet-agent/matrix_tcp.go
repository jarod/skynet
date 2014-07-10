package main

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"github.com/jarod/skynet/skynet"
	skc "github.com/jarod/skynet/skynet/client"
	skn "github.com/jarod/skynet/skynet/net"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	maxRetryDelay = time.Second * 32
)

type MatrixClient struct {
	conn *net.TCPConn
	cond *sync.Cond
}

func DialMatrix(raddr string) (mc *MatrixClient, err error) {
	addr, err := net.ResolveTCPAddr("tcp", raddr)
	if err != nil {
		return
	}
	mc, err = newMatrixClient()
	if err != nil {
		return
	}
	mc.connect(addr)
	return
}

func newMatrixClient() (mc *MatrixClient, err error) {
	mc = new(MatrixClient)
	mc.cond = sync.NewCond(&sync.Mutex{})
	return
}

func (mc *MatrixClient) connect(raddr *net.TCPAddr) {
	connDelay := time.Second
	for i := uint32(0); ; i++ {
		conn, err := net.DialTCP("tcp", nil, raddr)
		if err != nil {
			if connDelay < maxRetryDelay {
				connDelay *= 2
			}
			if i%4 == 0 {
				log.Printf("Failed to connect Matrix server %v, reconnect in %v", raddr, connDelay)
			}
			time.Sleep(connDelay)
		} else {
			log.Printf("Connected to Matrix[%s]", conn.RemoteAddr())
			go mc.onConnected(conn)
			break
		}
	}
}

func (m *MatrixClient) onConnected(conn *net.TCPConn) {
	m.cond.L.Lock()
	m.conn = conn
	m.cond.Broadcast()
	m.cond.L.Unlock()
	for {
		p, err := skn.ParsePacket(conn)
		if err != nil {
			if err != io.EOF {
				log.Printf("onConnected - %v\n", err)
			}
			break
		}
		m.dispatchMessage(p)
	}
	conn.Close()
	m.cond.L.Lock()
	m.conn = nil
	m.cond.L.Unlock()
	go m.connect(conn.RemoteAddr().(*net.TCPAddr))
}

func (mc *MatrixClient) Write(p *skn.Packet) {
	data := p.Encode()
	mc.cond.L.Lock()
	if mc.conn == nil {
		mc.cond.Wait()
	}
	io.Copy(mc.conn, bytes.NewReader(data))
	mc.cond.L.Unlock()
}

func (m *MatrixClient) dispatchMessage(p *skn.Packet) {
	switch p.Head {
	case 0x0000:
		m.onAppInfoUpdate(p)
	case 0x0001:
		m.onAppDisconnect(p)
	case 0x0002:
		m.execAgentCmd(p)
	default:
		tcpServer.BroadcastApps(p)
	}
}

func (m *MatrixClient) onAppInfoUpdate(p *skn.Packet) {
	info := new(skc.AppInfo)
	err := json.Unmarshal(p.Body, info)
	if err != nil {
		log.Println("onAppInfoUpdate - ", err)
		return
	}
	appInfos[info.Id] = info
}

func (m *MatrixClient) onAppDisconnect(p *skn.Packet) {
	id := new(skynet.Pstring)
	err := proto.Unmarshal(p.Body, id)
	if err != nil {
		log.Println("onAppDisconnect - ", err)
		return
	}
	delete(appInfos, id.GetValue())
}

func (m *MatrixClient) execAgentCmd(p *skn.Packet) {
	go func() {
		c := new(skynet.Pstring)
		proto.Unmarshal(p.Body, c)

		log.Println("execAgentCmd", c.GetValue())
		rawCmd := strings.Split(c.GetValue(), " ")
		cmd := exec.Command(rawCmd[0], rawCmd[1:]...)
		data, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		log.Println(string(data))
	}()
}
