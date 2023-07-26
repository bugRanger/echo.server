package listener

import (
	"context"
	"errors"
	"log"
	"net"
)

type UDPListener struct {
}

func (l *UDPListener) Listen(ctx context.Context, address string, handler PacketHandler) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if nil != err {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	go l.listen(conn, handler)

	<-ctx.Done()

	return nil
}

func (l *UDPListener) listen(conn *net.UDPConn, handler PacketHandler) {
	buf := make([]byte, 1024)
	for {
		count, addr, err := conn.ReadFromUDPAddrPort(buf)

		if count > 0 {
			packet := handler.Handle(buf[:count])

			_, err = conn.WriteToUDPAddrPort(packet, addr)
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
}
