package listener

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"syscall"
)

type TcpListener struct {
	handler ConnectionHandler
	closer  io.Closer
}

type ConnectionHandler interface {
	handle(conn net.Conn)
}

func NewTcpListener(address string, handler ConnectionHandler) (*TcpListener, error) {
	config := &net.ListenConfig{Control: reuseAddr}

	log.Println("Start requested")
	listener, err := config.Listen(context.Background(), "tcp", address)
	if err != nil {
		return nil, err
	}

	log.Println("Server is running on:", address)

	tcpListener := &TcpListener{
		handler,
		listener,
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

			return
		}

		go listener.handler.handle(conn)
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
