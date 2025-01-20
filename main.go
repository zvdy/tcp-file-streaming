package main

import (
	"github.com/zvdy/tcp-file-streaming/server"
)

func main() {
	server := &server.FileServer{Port: ":8080"}
	server.Start()
}
