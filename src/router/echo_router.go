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

	log.Println("Server UDP is running on:", address)

	tcpListener, err := listener.NewTcpListener(address, handler)
	if err != nil {
		return err
	}

	log.Println("Server TCP is running on:", address)

	e.tcp = tcpListener
	e.udp = udpListener

	return nil
}

func (e *EchoRouter) Close() {
	var err error
	log.Println("Stop requested")

	if e.udp != nil {
		err = e.udp.Close()
		if err != nil {
			log.Println("Failed close udp:", err.Error())
		}
	}

	if e.tcp != nil {
		err = e.tcp.Close()
		if err != nil {
			log.Println("Failed close udp:", err.Error())
		}
	}

	log.Println("Stopped successfully")
}
