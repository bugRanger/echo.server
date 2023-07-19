package main

import (
	"flag"
	"log"
	"main/src/listener"
	"time"
)

func main() {

	addressFlag := flag.String("l", "127.0.0.1:7", "a string as `127.0.0.1:7`")
	flag.Parse()

	listener, err := listener.NewTcpListener(*addressFlag, &listener.EchoHandler{})
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(30 * time.Second)

	err = listener.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
