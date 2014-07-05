package main

import (
	"encoding/json"
	skc "github.com/jarod/skynet/skynet/client"
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

func (h *HttpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//log.Println(req.URL.String())
	switch req.URL.Path {
	case "/agent/apps":
		h.findApps(w, req)
	default:
		log.Printf("http: no handler for request. %s %s\n", req.Method, req.RequestURI)
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
			enc.Encode(map[string]interface{}{
				"Code": 1,
				"Msg":  err.Error()})
			return
		}
		log.Printf("k:%s,v:%v,pattern:%s\n", k, v, pattern)
		if matched {
			infos = append(infos, v)
		}
	}

	err = enc.Encode(map[string]interface{}{
		"Code": 0,
		"Data": infos})
	if err != nil {
		log.Println("agent http: findApps - ", err)
		return
	}
}
