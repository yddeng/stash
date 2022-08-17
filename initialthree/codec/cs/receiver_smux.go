package cs

import (
	"github.com/sniperHW/flyfish/pkg/buffer"
	"io"
	"net"
)

func ReadMessage(conn net.Conn) ([]byte, error) {
	var err error
	var payload uint16

	hdr := make([]byte, SizeLen)
	if _, err = io.ReadFull(conn, hdr); err != nil {
		return nil, err
	}

	reader := buffer.NewReader(hdr)
	if payload, err = reader.CheckGetUint16(); err != nil {
		return nil, err
	}

	buf := make([]byte, payload+SizeLen)
	copy(buf, hdr)
	if _, err = io.ReadFull(conn, buf[SizeLen:]); err != nil {
		return nil, err
	}

	return buf, err
}

func Encode(msg *Message, encoder *Encoder) ([]byte, error) {
	buff := buffer.Get()
	if err := encoder.EnCode(msg, buff); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
