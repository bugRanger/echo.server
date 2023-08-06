package main

import (
	"flag"
	"log"
	"main/src/router"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addressFlag := flag.String("l", "127.0.0.1:7", "a string as `127.0.0.1:7`")
	flag.Parse()

	router := &router.EchoRouter{}
	router.Open(*addressFlag)
	defer router.Close()

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Signal: ", <-chSig)
}
