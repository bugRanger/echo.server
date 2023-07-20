package listener

import (
	"io"
	"net"
)

type EchoHandler struct {
}

func (handler *EchoHandler) handle(conn net.Conn) error {
	_, err := io.Copy(conn, conn)
	return err
}

func (handler *EchoHandler) handleByte(bytes []byte) (array []byte) {
	return bytes
}
