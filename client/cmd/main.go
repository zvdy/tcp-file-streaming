package main

import (
	"log"

	"github.com/zvdy/tcp-file-streaming/client"
)

func main() {
	err := client.SendFile(300000) // Specify the file size you want to send
	if err != nil {
		log.Fatal(err)
	}
}
