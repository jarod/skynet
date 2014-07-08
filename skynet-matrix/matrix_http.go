package main

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"github.com/jarod/skynet/skynet"
	skc "github.com/jarod/skynet/skynet/client"
	"github.com/jarod/skynet/skynet/net"
	"log"
	"net/http"
	"regexp"
)

type HttpServer struct {
}

func NewHttpServer() *HttpServer {
	return &HttpServer{}
}

func (h *HttpServer) ListenAndServe(addr string) {
	log.Println("Serve http on ", addr)
	log.Fatal(http.ListenAndServe(addr, h))
}

func (h *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.URL.String())
	switch r.URL.Path {
	case "/matrix/apps":
		h.findApps(w, r)
	case "/matrix/cmd":
		h.execAgentCmd(w, r)
	default:
		log.Printf("http: no handler for request. %s %s\n", r.Method, r.RequestURI)
	}
}

func (h *HttpServer) findApps(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Println("agent http: findApps - ", err)
		return
	}
	pattern := req.FormValue("pattern")
	var infos []*skc.AppInfo
	enc := json.NewEncoder(w)
	for k, v := range appInfos {
		matched, err := regexp.MatchString(pattern, k)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//log.Printf("k:%s,v:%v,pattern:%s\n", k, v, pattern)
		if matched {
			infos = append(infos, v)
		}
	}

	err = enc.Encode(map[string]interface{}{
		"Data": infos})
	if err != nil {
		log.Println("agent http: findApps - ", err)
		return
	}
}

/*
URL: /matrix/cmd?agent=agent_addr&cmd=
send command to agent and run on agent machine
@param agent: addr of agent server
@param cmd: shell command to execute
*/
func (h *HttpServer) execAgentCmd(w http.ResponseWriter, r *http.Request) {
	addr := r.FormValue("agent")
	mutex.Lock()
	agent := tcpServer.FindAgentByAddr(addr)
	mutex.Unlock()
	if agent == nil {
		log.Printf("No agent of addr=%s\n", addr)
		return
	}
	cmd := r.FormValue("cmd")
	log.Printf("exec agent command - agent=%s,cmd=%s,\n", addr, cmd)
	p, err := net.NewMessagePacket(0x0002, &skynet.Pstring{Value: proto.String(cmd)})
	if err != nil {
		log.Println(err)
		return
	}
	agent.Write(p)
}
