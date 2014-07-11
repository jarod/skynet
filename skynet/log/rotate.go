package log

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	file *os.File
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

/*
file: logging file name
e.g. /var/log/my_app, rotate file will be /var/log/my_app-yyyyMMddhhmmss
*/
func RegisterRotate(fn string) {
	if fn == "" {
		return
	}

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR1)
	go waitUsr1Signal(c)

	updateLogFile(fn)
}

func updateLogFile(fn string) {
	if file != nil {
		now := time.Now()
		rotated := fmt.Sprintf("%s-%d%02d%02d_%02d%02d%02d", fn, now.Year(), now.Month(), now.Day(), now.Hour(),
			now.Minute(), now.Second())
		file.Close()
		os.Rename(fn, rotated)
	}

	var err error
	file, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("updateLogFile(%s) - %s\n", fn, err))
		return
	}
	log.SetOutput(file)
}

func waitUsr1Signal(c chan os.Signal) {
	for {
		select {
		case <-c:
			updateLogFile(file.Name())
		}
	}
}
