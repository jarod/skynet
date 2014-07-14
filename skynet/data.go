package skynet

import (
	"fmt"
)

/**
 JSON:
 {
	Id: app id
	Gateway: {  // logic - gateway tcp socket, logic server only
	   Host:
	   Port:
	}
 }
*/
type AppInfo struct {
	Id     string
	Status uint8 // 0 connected, 100 disconnected
	Agent  string
}

func (a *AppInfo) String() string {
	return fmt.Sprintf("AppInfo{Id=%s,Agent=%s,Status=%d}", a.Id, a.Agent, a.Status)
}
