package server

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type FileServer struct {
	Port string
}

func (fs *FileServer) Start() {
	ln, err := net.Listen("tcp", fs.Port)
	fmt.Println("Server started on", fs.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("Error accepting connection:", err)
				continue
			}
			go fs.handleConnection(conn)
		}
	}()

	<-stop
	log.Println("Shutting down server...")
}

func (fs *FileServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := new(bytes.Buffer)
	for {
		var size int64
		err := binary.Read(conn, binary.LittleEndian, &size)
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed by client")
				return
			}
			log.Println("Error reading size:", err)
			return
		}
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed by client")
				return
			}
			log.Println("Error reading data:", err)
			return
		}

		// Generate a random unique ID string
		randomID := make([]byte, 16)
		_, err = rand.Read(randomID)
		if err != nil {
			log.Println("Error generating random ID:", err)
			return
		}
		fileName := fmt.Sprintf("%d_%x", size, randomID)

		// Log the received data
		fmt.Printf("Received %d bytes and saved to memory buffer with ID %s\n", n, fileName)

		// Notify the client that the server has finished processing
		_, err = conn.Write([]byte("done"))
		if err != nil {
			log.Println("Error notifying client:", err)
			return
		}

		// Close the connection
		return
	}
}
