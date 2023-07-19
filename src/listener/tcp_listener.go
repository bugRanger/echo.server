package listener

import (
	"context"
	"io"
	"log"
	"net"
	"syscall"
)

type TcpListener struct {
	listener net.Listener
	shutdown chan bool
}

func NewTcpListener(address string) *TcpListener {

	config := &net.ListenConfig{Control: reuseAddr}

	listener, err := config.Listen(context.Background(), "tcp", address)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Server is running on:", address)

	tcpListener := &TcpListener{
		listener,
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

	go handleConn(conn)
}

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	io.Copy(conn, conn)
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
