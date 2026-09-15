package main

import (
	"log"
	"server/client"
)

func main() {
	go client.Maintain()
	err := client.Listen()
	if err != nil {
		log.Printf("client error: %s\n", err.Error())
		return
	}
}
