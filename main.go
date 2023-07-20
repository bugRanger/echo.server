package main

import (
	"flag"
	"log"
	"main/src/router"
	"time"
)

func main() {

	addressFlag := flag.String("l", "127.0.0.1:7", "a string as `127.0.0.1:7`")
	flag.Parse()

	router := &router.EchoRouter{}

	err := router.Open(*addressFlag)
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(30 * time.Second)

	router.Close()
}
