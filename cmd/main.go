package main

import (
	"log"
	"os"

	"github.com/zvdy/tcp-file-streaming/server"
)

func main() {
	tcpPort := os.Getenv("FILE_SERVER_PORT")
	if tcpPort == "" {
		log.Fatal("FILE_SERVER_PORT environment variable is not set")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		log.Fatal("HTTP_PORT environment variable is not set")
	}

	fileServer := &server.FileServer{TCPPort: tcpPort, HTTPPort: httpPort}
	fileServer.Start()
}
