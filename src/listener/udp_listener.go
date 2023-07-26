package listener

import (
	"errors"
	"io"
	"log"
	"net"
)

type UdpListener struct {
	closer io.Closer
}

func NewUdpListener(address string, handler PacketHandler) (*UdpListener, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if nil != err {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
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
	}()

	udpListener := &UdpListener{
		conn,
	}

	return udpListener, nil
}

func (listener *UdpListener) Close() error {

	err := listener.closer.Close()
	if err != nil {
		return err
	}

	return nil
}
