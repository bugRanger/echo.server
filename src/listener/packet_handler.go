package listener

type PacketHandler interface {
	Handle(bytes []byte) (array []byte)
}
