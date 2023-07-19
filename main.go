package main

import (
	"flag"
	"main/src/listener"
	"time"
)

func main() {

	addressFlag := flag.String("l", "127.0.0.1:7", "a string as `127.0.0.1:7`")
	flag.Parse()

	listener := listener.NewTcpListener(*addressFlag)
	time.Sleep(30 * time.Second)
	listener.Close()
}
