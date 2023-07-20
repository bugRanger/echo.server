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
	handler ConnectionHandler
	closer  io.Closer
	conn    *sync.Map
}

type ConnectionHandler interface {
	handle(conn net.Conn) error
	handleByte(bytes []byte) (array []byte)
}

func NewTcpListener(address string, handler ConnectionHandler) (*TcpListener, error) {
	config := &net.ListenConfig{Control: reuseAddr}

	log.Println("Start requested")
	listener, err := config.Listen(context.Background(), "tcp", address)
	if err != nil {
		return nil, err
	}

	log.Println("Server is running on:", address)

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
	log.Println("Stop requested")

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

	log.Println("Stopped successfully")
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

func handleConnection(connId uuid.UUID, conn net.Conn, connSync *sync.Map, handler ConnectionHandler) {
	defer func() {
		connSync.Delete(connId)
		_ = conn.Close()
	}()

	connSync.Store(connId, conn)

	err := handler.handle(conn)
	if err != nil {
		if !errors.Is(err, net.ErrClosed) {
			log.Println("Failed to handle connection:", err.Error())
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
