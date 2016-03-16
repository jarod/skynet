package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jarod/skynet/skynet"
	skn "github.com/jarod/skynet/skynet/net"
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
	r := bytes.NewReader(p.Encode())
	mc.cond.L.Lock()
	if mc.conn == nil {
		mc.cond.Wait()
	}
	c := mc.conn
	mc.cond.L.Unlock()
	io.Copy(c, r)
}

func (mc *MatrixClient) Ping() {
	p := skn.NewEmptyPacket(0x0022)
	mc.Write(p)
}

func (m *MatrixClient) dispatchMessage(p *skn.Packet) {
	switch skynet.SkynetMsg(p.Head) {
	case skynet.SkynetMsg_SM_APP_INFO:
		m.onAppInfoUpdate(p)
	case skynet.SkynetMsg_SM_APP_DISCONNECTED:
		m.onAppDisconnect(p)
	case skynet.SkynetMsg_SM_AGENT_EXECUTE_CMD:
		m.execAgentCmd(p)
	case skynet.SkynetMsg_SM_AGENT_FIND_APPS:
		m.findApps(p)
	case skynet.SkynetMsg_SM_SEND_TO_APP:
		m.sendToApp(p)
	default:
		tcpServer.BroadcastApps(p)
	}
}

func (m *MatrixClient) onAppInfoUpdate(p *skn.Packet) {
	info := new(skynet.AppInfo)
	err := json.Unmarshal(p.Body, info)
	if err != nil {
		log.Println("onAppInfoUpdate - ", err)
		return
	}
	appInfos[info.Id] = info

	tcpServer.BroadcastApps(p)
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
		err := proto.Unmarshal(p.Body, c)
		if err != nil {
			log.Println("execAgentCmd - ", err)
			return
		}
		log.Println("execAgentCmd", c.GetValue())
		rawCmd := strings.Split(c.GetValue(), " ")
		cmd := exec.Command(rawCmd[0], rawCmd[1:]...)
		data, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(data))
	}()
}

func (m *MatrixClient) findApps(p *skn.Packet) {
	pattern := string(p.Body)
	infos, err := FindApps(pattern)
	if err != nil {
		log.Println("findApps - ", err)
		return
	}
	if len(infos) < 1 {
		return
	}
	p.Body, err = json.Marshal(infos)
	if err != nil {
		log.Println("findApps - ", err)
		return
	}
	m.Write(p)
}

func (m *MatrixClient) sendToApp(p *skn.Packet) {
	tcpServer.SendToApp(p)
}
