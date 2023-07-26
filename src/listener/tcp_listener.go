package listener

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type TCPListener struct {
	connections sync.WaitGroup
}

func (l *TCPListener) Listen(ctx context.Context, address string, handler PacketHandler) error {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if nil != err {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		l.connections.Wait()
	}()

	listener, err := net.ListenTCP("tcp", addr)
	if nil != err {
		return err
	}
	defer listener.Close()

	go l.listen(ctx, listener, handler)

	<-ctx.Done()

	return nil
}

func (l *TCPListener) listen(ctx context.Context, listener net.Listener, handler PacketHandler) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			if e, ok := err.(*net.OpError); ok && e.Temporary() {
				continue
			}

			if !errors.Is(err, net.ErrClosed) {
				log.Println("Failed to accept connection:", err.Error())
			}

			break
		}

		go l.handleConnection(ctx, conn, handler)
	}
}

func (l *TCPListener) handleConnection(ctx context.Context, conn net.Conn, handler PacketHandler) {
	ctx, cancel := context.WithCancel(ctx)

	l.connections.Add(1)
	defer func() {
		_ = conn.Close()
		l.connections.Done()
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
				if errors.Is(err, net.ErrClosed) {
					return
				}

				if errors.Is(err, io.EOF) {
					return
				}

				log.Println("Failed to read connection:", err.Error())
				return
			}
		}
	}()

	<-ctx.Done()
}
