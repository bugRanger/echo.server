package listener

import (
	"io"
	"net"
)

type EchoHandler struct {
}

func (handler *EchoHandler) handle(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	io.Copy(conn, conn)
}
