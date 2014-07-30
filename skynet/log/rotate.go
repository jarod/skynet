package log

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	SIGUSR1 = syscall.Signal(0xa)
)

var (
	file     *os.File
	filename string
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

/*SetLogFile .
file: logging file name
*/
func SetLogFile(fn string) {
	if fn == "" {
		return
	}
	filename = fn

	c := make(chan os.Signal, 1)
	signal.Notify(c, SIGUSR1)
	go waitUsr1Signal(c)
	openLogFile(fn)
}

func openLogFile(fn string) {
	if file != nil {
		file.Close()
	}

	var err error
	file, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("openLogFile(%s) - %s\n", fn, err))
		return
	}
	log.SetOutput(file)
}

func waitUsr1Signal(c chan os.Signal) {
	for {
		select {
		case <-c:
			openLogFile(filename)
		}
	}
}
