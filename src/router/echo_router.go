package router

import (
	"context"
	"log"
	"main/src/listener"
	"sync"
)

type EchoRouter struct {
	channels sync.WaitGroup
	cancel   context.CancelFunc
}

func (e *EchoRouter) Open(address string) error {
	log.Println("Start requested")

	handler := &EchoHandler{}

	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel

	e.channels.Add(2)
	go e.listenUDP(ctx, address, handler)
	go e.listenTCP(ctx, address, handler)

	return nil
}

func (e *EchoRouter) Close() {
	log.Println("Stop requested")

	e.cancel()
	e.channels.Wait()

	log.Println("Stopped successfully")
}

func (e *EchoRouter) listenTCP(ctx context.Context, address string, handler listener.PacketHandler) {
	defer e.channels.Done()

	log.Println("Server TCP is running on:", address)

	var listener listener.TCPListener
	err := listener.Listen(ctx, address, handler)

	log.Println("Server TCP stopped:", address)
	if err != nil {
		log.Println("Server TCP failure:", err)
	}
}

func (e *EchoRouter) listenUDP(ctx context.Context, address string, handler listener.PacketHandler) {
	defer e.channels.Done()

	log.Println("Server UDP is running on:", address)

	var listener listener.UDPListener
	err := listener.Listen(ctx, address, handler)

	log.Println("Server UDP stopped:", address)
	if err != nil {
		log.Println("Server UDP failure:", err)
	}
}
