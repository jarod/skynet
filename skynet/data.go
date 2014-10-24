package skynet

import (
	"fmt"
)

type AppInfo struct {
	Id      string
	Status  uint8 // 0 connected, 50 inactive, 100 disconnected
	Agent   string
	Gateway *struct {
		Host string
		Port uint16
	}
}

func (a *AppInfo) String() string {
	return fmt.Sprintf("AppInfo{Id=%s,Agent=%s,Status=%d}", a.Id, a.Agent, a.Status)
}
