package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"syscall"
)

func main() {

	addressFlag := flag.String("l", "127.0.0.1:7", "a string as `127.0.0.1:7`")
	flag.Parse()

	address := *addressFlag
	config := &net.ListenConfig{Control: reuseAddr}

	server, err := config.Listen(context.Background(), "tcp", address)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		log.Println("Server stopping")
		server.Close()
		log.Println("Server stopped")
	}()

	log.Println("Server is running on:", address)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Failed to accept conn", err)
			continue
		}

		go func(conn net.Conn) {
			defer func() {
				conn.Close()
			}()
			io.Copy(conn, conn)
		}(conn)
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
