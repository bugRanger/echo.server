package listener

import (
	"context"
	"log"
	"net"
	"syscall"
)

type TcpListener struct {
	listener net.Listener
	handler  ConnectionHandler
	shutdown chan bool
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
		make(chan bool),
	}

	go tcpListener.waitTerminate()
	go tcpListener.listen()

	return tcpListener
}

func (listener *TcpListener) Close() {
	log.Println("Stop requested")
	listener.shutdown <- true

	<-listener.shutdown

	log.Println("Stopped successfully")
}

func (listener *TcpListener) waitTerminate() {

	<-listener.shutdown
	log.Println("Shutting down...")
	listener.listener.Close()
	listener.shutdown <- true
	return
}

func (listener *TcpListener) listen() {

	conn, err := listener.listener.Accept()
	if err != nil {
		log.Println("Failed to accept connection:", err.Error())
	}

	go listener.handler.handle(conn)
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
