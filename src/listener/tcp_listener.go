package listener

import (
	"context"
	"errors"
	"log"
	"net"
	"syscall"
)

type TcpListener struct {
	listener net.Listener
	handler  ConnectionHandler
}

type ConnectionHandler interface {
	handle(conn net.Conn)
}

func NewTcpListener(address string, handler ConnectionHandler) *TcpListener {

	config := &net.ListenConfig{Control: reuseAddr}

	listener, err := config.Listen(context.Background(), "tcp", address)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Server is running on:", address)

	tcpListener := &TcpListener{
		listener,
		handler,
	}

	go tcpListener.listen()

	return tcpListener
}

func (listener *TcpListener) Close() {
	log.Println("Stop requested")

	listener.listener.Close()

	log.Println("Stopped successfully")
}

func (listener *TcpListener) listen() {
	for {
		conn, err := listener.listener.Accept()
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
