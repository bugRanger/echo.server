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

func NewUdpListener(address string, handler ConnectionHandler) (*UdpListener, error) {
	log.Println("Start requested")

	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return nil, err
	}

	log.Println("Server is running on:", address)

	go func() {
		buf := make([]byte, 1024)
		for {
			count, addr, err := conn.ReadFrom(buf)

			if count > 0 {
				packet := handler.handleByte(buf[:count])

				_, err = conn.WriteTo(packet, addr)
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
	log.Println("Stop requested")

	err := listener.closer.Close()
	if err != nil {
		return err
	}

	log.Println("Stopped successfully")
	return nil
}
