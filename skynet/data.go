package skynet

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
