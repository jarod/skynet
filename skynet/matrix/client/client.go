package client

import (
	"bytes"
	snet "github.com/jarod/skynet/skynet/net"
	"io"
	"log"
	"net"
	"time"
)

type MatrixClient struct {
	raddr *net.TCPAddr
	conn  *net.TCPConn

	connCh chan *net.TCPAddr
}

func Dial(serverAddr string) (mc *MatrixClient, err error) {
	mc, err = newMatrixClient(serverAddr)
	if err != nil {
		return
	}
	doneCh := make(chan bool)
	go mc.connect(doneCh)
	<-doneCh
	return
}

func newMatrixClient(serverAddr string) (mc *MatrixClient, err error) {
	addr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		return
	}

	mc = new(MatrixClient)
	mc.raddr = addr
	mc.connCh = make(chan *net.TCPAddr, 1)
	mc.connCh <- addr
	return
}

func (mc *MatrixClient) connect(doneCh chan bool) {
	connDelay := time.Duration(1)
	first := true
	for {
		addr := <-mc.connCh
		time.AfterFunc(connDelay*time.Second, func() {
			conn, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				if connDelay < 32 {
					connDelay *= 2
				}
				log.Printf("Failed to connect Matrix server %v, reconnect in %d seconds", addr, connDelay)
				mc.connCh <- addr
			} else {
				connDelay = 1
				mc.conn = conn
				log.Printf("Connected to Matrix[%s]", conn.RemoteAddr())
				if first {
					first = false
					doneCh <- true
				}
			}
		})
	}
}

func (mc *MatrixClient) Read() (p *snet.Packet, err error) {
	if mc.conn == nil {
		err = io.EOF
		return
	}
	p, err = snet.ParsePacket(mc.conn)
	if err != nil {
		mc.conn.Close()
		mc.connCh <- mc.raddr
		log.Printf("Read %v\n", err)
	}
	return
}

func (mc *MatrixClient) Write(p *snet.Packet) {
	data := p.Encode()
	io.Copy(mc.conn, bytes.NewReader(data))
}
