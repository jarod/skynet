package net

import (
	proto "code.google.com/p/goprotobuf/proto"
	"fmt"
	"io"
)

type Packet struct {
	Head uint16
	Body []byte
}

func ParsePacket(r io.Reader) (p *Packet, err error) {
	headerbuf := make([]byte, 5)
	_, err = io.ReadFull(r, headerbuf)
	if err != nil {
		return
	}

	blen := uint32(headerbuf[0]) << 16
	blen += uint32(headerbuf[1]) << 8
	blen += uint32(headerbuf[2])

	head := uint16(headerbuf[3]) << 8
	head += uint16(headerbuf[4])
	body := make([]byte, blen)
	_, err = io.ReadFull(r, body)
	if err != nil {
		return
	}
	p = NewPacket(head, body)
	return
}

func NewPacket(head uint16, body []byte) *Packet {
	p := new(Packet)
	p.Head = head
	p.Body = body
	return p
}

func NewMessagePacket(head uint16, pb proto.Message) (*Packet, error) {
	data, err := proto.Marshal(pb)
	if err != nil {
		return nil, err
	}
	p := NewPacket(head, data)
	return p, nil
}

func (p *Packet) String() string {
	return fmt.Sprintf("head=%02X,len=%d", p.Head, len(p.Body))
}

/* TODO performance */
func (p *Packet) Encode() []byte {
	l := len(p.Body)
	data := make([]byte, 5)
	data[0] = byte(0xFF & (l >> 16))
	data[1] = byte(0xFF & (l >> 8))
	data[2] = byte(0xFF & l)
	data[3] = byte(0xFF & (p.Head >> 8))
	data[4] = byte(0xFF & p.Head)
	return append(data, p.Body...)
}
