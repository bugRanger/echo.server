package listener

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type TcpListener struct {
	connections sync.WaitGroup
	closer      io.Closer
	cancel      context.CancelFunc
}

func NewTcpListener(address string, handler PacketHandler) (*TcpListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if nil != err {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if nil != err {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	tcpListener := &TcpListener{
		sync.WaitGroup{},
		listener,
		cancel,
	}

	go listen(ctx, listener, &tcpListener.connections, handler)

	return tcpListener, nil
}

func (listener *TcpListener) Close() error {
	err := listener.closer.Close()
	if err != nil {
		return err
	}

	listener.cancel()
	listener.connections.Wait()

	return nil
}

func listen(ctx context.Context, listenerBase net.Listener, connections *sync.WaitGroup, handler PacketHandler) {
	for {
		conn, err := listenerBase.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Println("Failed to accept connection:", err.Error())
			}

			continue
		}

		go handleConnection(ctx, conn, connections, handler)
	}
}

func handleConnection(ctx context.Context, conn net.Conn, connections *sync.WaitGroup, handler PacketHandler) {
	ctx, cancel := context.WithCancel(ctx)

	connections.Add(1)
	defer func() {
		_ = conn.Close()
		connections.Done()
	}()

	go func() {
		defer cancel()

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
	}()

	<-ctx.Done()
}
