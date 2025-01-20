package client

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func SendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	defer conn.Close()

	binary.Write(conn, binary.LittleEndian, int64(size))
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		return err
	}

	fmt.Printf("Sent %d bytes\n", n)

	// Wait for the server to notify that it has finished processing
	buf := make([]byte, 4)
	_, err = conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Server closed the connection")
			return nil
		}
		return err
	}

	if string(buf) == "done" {
		fmt.Println("Server finished processing")
	}

	return nil
}
