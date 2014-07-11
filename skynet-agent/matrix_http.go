package main

import (
	"encoding/json"
	"github.com/jarod/skynet/skynet"
	"log"
	"net/http"
)

type MatrixHttpClient struct {
	urlBase string
}

func NewMatrixHttpClient(urlBase string) *MatrixHttpClient {
	return &MatrixHttpClient{urlBase: urlBase}
}

func (m *MatrixHttpClient) FindApps(pattern string) []*skynet.AppInfo {
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
