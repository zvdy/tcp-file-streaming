package tcp

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

type TCPServer struct {
	Port string
}

func (ts *TCPServer) Start() {
	ln, err := net.Listen("tcp", ts.Port)
	fmt.Println("TCP server started on", ts.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
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

		randomID := make([]byte, 16)
		_, err = rand.Read(randomID)
		if err != nil {
			log.Println("Error generating random ID:", err)
			return
		}
		fileName := fmt.Sprintf("%d_%x", size, randomID)

		fmt.Printf("Received %d bytes and saved to memory buffer with ID %s\n", n, fileName)

		_, err = conn.Write([]byte("done"))
		if err != nil {
			log.Println("Error notifying client:", err)
			return
		}
	}
}

func SendToTCPServer(data []byte) error {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	defer conn.Close()

	dataSize := int64(len(data))
	err = binary.Write(conn, binary.LittleEndian, dataSize)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	response := make([]byte, 4)
	_, err = conn.Read(response)
	if err != nil {
		return err
	}

	if string(response) != "done" {
		return fmt.Errorf("unexpected response from server: %s", string(response))
	}

	return nil
}
