package main

import (
	"bytes"
	"encoding/json"
	"github.com/jarod/skynet/skynet"
	skn "github.com/jarod/skynet/skynet/net"
	"io"
	"log"
	"net"
	"regexp"
)

var (
	appInfos map[string]*skynet.AppInfo
)

func init() {
	appInfos = make(map[string]*skynet.AppInfo)
}

func FindApps(pattern string) (infos []*skynet.AppInfo, err error) {
	for k, v := range appInfos {
		matched := false
		matched, err = regexp.MatchString(pattern, k)
		if err != nil {
			return
		}
		//log.Printf("k:%s,v:%v,pattern:%s\n", k, v, pattern)
		if matched {
			infos = append(infos, v)
		}
	}
	return
}

type App struct {
	id   string
	conn *net.TCPConn
}

func NewApp(conn *net.TCPConn) (ac *App) {
	ac = new(App)
	ac.conn = conn
	return
}

func (ac *App) Write(p *skn.Packet) {
	data := p.Encode()
	io.Copy(ac.conn, bytes.NewReader(data))
}

func (ac *App) dispatchAppPacket(p *skn.Packet) {
	switch skynet.SkynetMsg(p.Head) {
	case skynet.SkynetMsg_SM_APP_INFO:
		ac.updateAppInfo(p)
	case skynet.SkynetMsg_SM_SEND_TO_APP:
		tcpServer.SendToApp(p)
	default:
		matrixClient.Write(p)
	}
	//log.Printf("dispatchAppPacket %v\n", p)
}

func (a *App) updateAppInfo(p *skn.Packet) {
	info := new(skynet.AppInfo)
	err := json.Unmarshal(p.Body, info)
	if err != nil {
		log.Println("app: updateAppInfo - ", err)
		return
	}
	a.id = info.Id
	tcpServer.appIdMap[a.id] = a
	matrixClient.Write(p)
}
