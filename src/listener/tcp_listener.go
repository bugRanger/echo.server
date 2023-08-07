package listener

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"time"
)

type TCPListener struct {
}

func (l *TCPListener) Listen(ctx context.Context, address string, handler PacketHandler) error {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	chanError := make(chan error, 1)
	go func() {
		defer cancel()

		err := l.listen(ctx, listener, handler)
		if err != nil {
			chanError <- err
		}

		close(chanError)
	}()

	<-ctx.Done()
	listener.Close()

	err, ok := <-chanError
	if ok && err != nil {
		return err
	}

	return nil
}

func (l *TCPListener) listen(ctx context.Context, listener net.Listener, handler PacketHandler) error {
	var tempDelay time.Duration

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}

			if e, ok := err.(*net.OpError); ok && e.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				time.Sleep(tempDelay)
				continue
			}

			return err
		}

		go l.handleConnection(ctx, conn, handler)
	}
}

func (l *TCPListener) handleConnection(ctx context.Context, conn net.Conn, handler PacketHandler) {
	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		_ = conn.Close()
	}()

	go func() {
		defer cancel()

		buf := make([]byte, 256)
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
