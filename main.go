package main

import (
	"log"
	"sonic/channel"
)

func main() {
	channel, err := channel.NewInjestChannel("[::1]:1491")
	if err != nil {
		log.Fatalf("Could not make new channel: %s", err.Error())
	}

	channel.Verbose = true

	channel.Start("SecretPassword")
	log.Println("Started channel")

	err = channel.Push("messages", "user:artem", "conversation:123", "Hello world!")
	if err != nil {
		log.Println(err.Error())
	}

	err = channel.Ping()
	if err != nil {
		log.Println(err.Error())
	}

	channel.Quit()
	log.Println("Stopped channel")
}
