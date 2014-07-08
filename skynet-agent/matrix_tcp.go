package main

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"github.com/jarod/skynet/skynet"
	skc "github.com/jarod/skynet/skynet/client"
	skmc "github.com/jarod/skynet/skynet/matrix/client"
	skn "github.com/jarod/skynet/skynet/net"
	"log"
	"os/exec"
	"strings"
)

type MatrixClient struct {
	mc *skmc.MatrixClient
}

func DialMatrix(addr string) (*MatrixClient, error) {
	mc, err := skmc.Dial(addr)
	if err != nil {
		return nil, err
	}
	c := &MatrixClient{mc: mc}
	go c.readMatrix()
	return c, nil
}

func (m *MatrixClient) readMatrix() {
	for {
		p, err := m.mc.Read()
		if err != nil {
			log.Printf("readMatrix err=%v", err)
			break
		}
		m.dispatchMessage(p)
	}
}

func (m *MatrixClient) Write(p *skn.Packet) {
	m.mc.Write(p)
}

func (m *MatrixClient) dispatchMessage(p *skn.Packet) {
	switch p.Head {
	case 0x0000:
		m.onMatrixAppInfoUpdate(p)
	case 0x0001:
		m.onMatrixAppDisconnect(p)
	case 0x0002:
		m.execAgentCmd(p)
	default:
		tcpServer.BroadcastApps(p)
	}
}

func (m *MatrixClient) onMatrixAppInfoUpdate(p *skn.Packet) {
	info := new(skc.AppInfo)
	err := json.Unmarshal(p.Body, info)
	if err != nil {
		log.Println("onMatrixAppInfoUpdate - ", err)
		return
	}
	appInfos[info.Id] = info
}

func (m *MatrixClient) onMatrixAppDisconnect(p *skn.Packet) {
	id := new(skynet.Pstring)
	err := proto.Unmarshal(p.Body, id)
	if err != nil {
		log.Println("onMatrixAppDisconnect - ", err)
		return
	}
	delete(appInfos, id.GetValue())
}

func (m *MatrixClient) execAgentCmd(p *skn.Packet) {
	go func() {
		c := new(skynet.Pstring)
		proto.Unmarshal(p.Body, c)

		log.Println("exec cmd=", c.GetValue())
		rawCmd := strings.Split(c.GetValue(), " ")
		cmd := exec.Command(rawCmd[0], rawCmd[1:]...)
		data, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		log.Println(string(data))
	}()
}
