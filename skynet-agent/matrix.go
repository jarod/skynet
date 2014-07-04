package main

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"github.com/jarod/skynet/skynet"
	skc "github.com/jarod/skynet/skynet/client"
	skmc "github.com/jarod/skynet/skynet/matrix/client"
	"github.com/jarod/skynet/skynet/net"
	"log"
	"os/exec"
	"strings"
	"time"
)

func readMatrix(mc *skmc.MatrixClient) {
	matrixClient = mc
	for {
		p, err := mc.Read()
		if err != nil {
			log.Printf("readMatrix err=%v", err)
			time.Sleep(16 * time.Second)
			continue
		}
		dispatchMatrixMessage(p)
	}
}

func dispatchMatrixMessage(p *net.Packet) {
	switch p.Head {
	case 0x0000:
		onMatrixAppInfoUpdate(p)
	case 0x0001:
		onMatrixAppDisconnect(p)
	case 0x0002:
		execAgentCmd(p)
	default:
		broadcastClients(p)
	}
}

func onMatrixAppInfoUpdate(p *net.Packet) {
	info := new(skc.AppInfo)
	err := json.Unmarshal(p.Body, info)
	if err != nil {
		log.Println("onMatrixAppInfoUpdate - ", err)
		return
	}
	appInfos[info.Id] = info
}

func onMatrixAppDisconnect(p *net.Packet) {
	id := new(skynet.Pstring)
	err := proto.Unmarshal(p.Body, id)
	if err != nil {
		log.Println("onMatrixAppDisconnect - ", err)
		return
	}
	delete(appInfos, id.GetValue())
}

func broadcastClients(p *net.Packet) {
	for _, v := range connMap {
		v.Write(p)
	}
}

func execAgentCmd(p *net.Packet) {
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
