package router

import (
	"log"
	"main/src/listener"
)

type EchoRouter struct {
	tcp *listener.TcpListener
	udp *listener.UdpListener
}

func (e *EchoRouter) Open(address string) error {
	log.Println("Start requested")
	handler := &EchoHandler{}

	udpListener, err := listener.NewUdpListener(address, handler)
	if err != nil {
		return err
	}
	e.udp = udpListener

	log.Println("Server UDP is running on:", address)

	tcpListener, err := listener.NewTcpListener(address, handler)
	if err != nil {
		return err
	}

	e.tcp = tcpListener

	log.Println("Server TCP is running on:", address)

	return nil
}

func (e *EchoRouter) Close() {
	var err error
	log.Println("Stop requested")

	if e.udp != nil {
		err = e.udp.Close()
		if err != nil {
			log.Println("Failed close server UDP:", err.Error())
		}
	}

	if e.tcp != nil {
		err = e.tcp.Close()
		if err != nil {
			log.Println("Failed close server TCP:", err.Error())
		}
	}

	log.Println("Stopped successfully")
}
