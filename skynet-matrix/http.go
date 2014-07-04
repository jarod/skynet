package main

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"github.com/jarod/skynet/skynet"
	"github.com/jarod/skynet/skynet/net"
	"log"
	"net/http"
	"strings"
)

func init() {
	http.HandleFunc("/clients", getClientList)
	http.HandleFunc("/cmd", execAgentCmd)
}

type MatrixHttpServer struct {
}

func NewMatrixHttpServer() *MatrixHttpServer {
	return &MatrixHttpServer{}
}

func (m *MatrixHttpServer) Startup(addr string) {
	log.Println("Serve http on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

type client struct {
	Id string `json:"id"`
	Ip string `json:"ip"`
}

/*
URL: /clients
get running client ids
@return [{id:ip},{...}]
*/
func getClientList(w http.ResponseWriter, r *http.Request) {
	clients := make([]*client, 0, len(apps))
	for k, v := range apps {
		parts := strings.Split(v.Agent, ":")
		if len(parts) != 2 {
			continue
		}
		clients = append(clients, &client{Id: k, Ip: parts[0]})
	}

	data, err := json.Marshal(clients)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(data)
}

/*
URL: /cmd?ip=agent_ip&cmd=
send command to agent and run on agent machine
@param ip: ip of agent server
@param cmd: shell command to execute
*/
func execAgentCmd(w http.ResponseWriter, r *http.Request) {
	ip := r.FormValue("ip")
	mutex.Lock()
	agent := findAgentByIP(ip)
	mutex.Unlock()
	if agent == nil {
		log.Printf("No agent of ip=%s\n", ip)
		return
	}
	cmd := r.FormValue("cmd")
	log.Printf("exec agent command - agent=%s,cmd=%s,\n", agent.RemoteIp(), cmd)
	p, err := net.NewMessagePacket(0x0002, &skynet.Pstring{Value: proto.String(cmd)})
	if err != nil {
		log.Println(err)
		return
	}
	agent.Write(p)
}
