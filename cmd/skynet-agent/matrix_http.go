package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jarod/skynet/skynet"
)

type MatrixHTTPClient struct {
	urlBase string
}

func NewMatrixHTTPClient(urlBase string) *MatrixHTTPClient {
	return &MatrixHTTPClient{urlBase: urlBase}
}

func (m *MatrixHTTPClient) FindApps(pattern string) []*skynet.AppInfo {
	resp, err := http.Get(m.urlBase + "matrix/apps?pattern=" + pattern)
	if err != nil {
		log.Println("FindApps - ", err)
		return nil
	}
	if resp.StatusCode != 200 {
		log.Printf("FindApps(%s) - %d %s\n", pattern, resp.StatusCode, resp.Status)
		return nil
	}
	ret := struct {
		Data []*skynet.AppInfo
	}{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ret)
	if err != nil {
		log.Println("FindApps - ", err)
		return nil
	}
	return ret.Data
}
