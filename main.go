package main

import (
	"github.com/zvdy/tcp-file-streaming/server"
)

func main() {
	server := &server.FileServer{}
	server.Start()
}
