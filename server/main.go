package main

import (
	"log"
	"os"
	"server/admin"
	"server/client"
	"server/config"
)

func init() {
	cfg, err := config.ReadConfig("config.json")
	if err != nil {
		log.Printf("read config error: %s\r\n", err.Error())
		os.Exit(-1)
	}

	go client.Maintain()
	go admin.ListenWeb(cfg)
}

func main() {
	err := client.Listen()
	if err != nil {
		log.Printf("client error: %s\n", err.Error())
		return
	}
}
