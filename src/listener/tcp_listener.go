package listener

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"syscall"

	"github.com/google/uuid"
)

type TcpListener struct {
	handler PacketHandler
	closer  io.Closer
	conn    *sync.Map
}

func NewTcpListener(address string, handler PacketHandler) (*TcpListener, error) {
	config := &net.ListenConfig{Control: reuseAddr}
	listener, err := config.Listen(context.Background(), "tcp", address)
	if err != nil {
		return nil, err
	}

	var conn sync.Map

	tcpListener := &TcpListener{
		handler,
		listener,
		&conn,
	}

	go tcpListener.listen(listener)

	return tcpListener, nil
}

func (listener *TcpListener) Close() error {
	err := listener.closer.Close()
	if err != nil {
		return err
	}

	listener.conn.Range(func(key, value interface{}) bool {
		val, ok := value.(io.Closer)
		if ok {
			val.Close()
		}

		return true
	})

	return nil
}

func (listener *TcpListener) listen(listenerBase net.Listener) {
	for {
		conn, err := listenerBase.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Println("Failed to accept connection:", err.Error())
			}

			continue
		}

		go handleConnection(uuid.New(), conn, listener.conn, listener.handler)
	}
}

func handleConnection(connId uuid.UUID, conn net.Conn, connSync *sync.Map, handler PacketHandler) {
	defer func() {
		connSync.Delete(connId)
		_ = conn.Close()
	}()

	connSync.Store(connId, conn)

	buf := make([]byte, 1024)
	for {
		count, err := conn.Read(buf)

		if count > 0 {
			packet := handler.Handle(buf[:count])

			_, err = conn.Write(packet)
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					log.Println("Failed to write connection:", err.Error())
				}

				return
			}
		}

		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Println("Failed to read connection:", err.Error())
			}

			return
		}
	}
}

func reuseAddr(network, address string, conn syscall.RawConn) error {
	var errorControl, errorReuseAddr error

	errorControl = conn.Control(func(descriptor uintptr) {
		errorReuseAddr = syscall.SetsockoptInt(syscall.Handle(descriptor), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	})

	if errorControl != nil {
		return errorControl
	}

	return errorReuseAddr
}
