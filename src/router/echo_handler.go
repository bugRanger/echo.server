package router

type EchoHandler struct {
}

func (handler *EchoHandler) Handle(bytes []byte) (array []byte) {
	return bytes
}
