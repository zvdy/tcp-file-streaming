package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zvdy/tcp-file-streaming/internal/consul"
	"github.com/zvdy/tcp-file-streaming/internal/http"
	"github.com/zvdy/tcp-file-streaming/internal/tcp"
)

type FileServer struct {
	TCPPort  string
	HTTPPort string
}

func (fs *FileServer) Start() {
	err := consul.RegisterService(fs.TCPPort, fs.HTTPPort)
	if err != nil {
		log.Fatal("Failed to register service:", err)
	}

	tcpServer := &tcp.TCPServer{Port: fs.TCPPort}
	go tcpServer.Start()

	httpServer := &http.HTTPServer{Port: fs.HTTPPort}
	go httpServer.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server...")
}
