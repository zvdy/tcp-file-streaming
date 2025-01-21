package server

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/hashicorp/consul/api"
)

type FileServer struct {
	TCPPort  string
	HTTPPort string
}

func (fs *FileServer) Start() {
	err := fs.registerService()
	if err != nil {
		log.Fatal("Failed to register service:", err)
	}

	// Start the TCP server
	go fs.startTCPServer()

	// Start the HTTP server for health checks and streaming
	go fs.startHTTPServer()

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server...")
}

func (fs *FileServer) startTCPServer() {
	ln, err := net.Listen("tcp", fs.TCPPort)
	fmt.Println("TCP server started on", fs.TCPPort)
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
		go fs.handleConnection(conn)
	}
}

func (fs *FileServer) startHTTPServer() {
	http.HandleFunc("/health", fs.HealthCheckHandler)
	http.HandleFunc("/stream", fs.StreamHandler)
	log.Println("HTTP server started on port", fs.HTTPPort)
	if err := http.ListenAndServe(fs.HTTPPort, nil); err != nil {
		log.Fatal("HTTP server error:", err)
	}
}

func (fs *FileServer) StreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Log request headers
	log.Println("Stream request headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("%s: %s\n", name, value)
		}
	}

	// Read the size parameter from the request body
	sizeStr := r.URL.Query().Get("size")
	if sizeStr == "" {
		http.Error(w, "Size parameter is missing", http.StatusBadRequest)
		return
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		http.Error(w, "Invalid size parameter", http.StatusBadRequest)
		return
	}

	// Generate random file data
	file := make([]byte, size)
	_, err = io.ReadFull(rand.Reader, file)
	if err != nil {
		http.Error(w, "Failed to generate file data", http.StatusInternalServerError)
		return
	}

	// Send the file data to the TCP server
	err = fs.sendToTCPServer(file)
	if err != nil {
		http.Error(w, "Failed to send data to TCP server", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data processed successfully"))
}

func (fs *FileServer) sendToTCPServer(data []byte) error {
	conn, err := net.Dial("tcp", "localhost"+fs.TCPPort)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send the size of the data
	dataSize := int64(len(data))
	err = binary.Write(conn, binary.LittleEndian, dataSize)
	if err != nil {
		return err
	}

	// Send the actual data
	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	// Wait for the server to acknowledge the data
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
	}
}

func (fs *FileServer) registerService() error {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	// Get the container's IP address
	ip, err := getContainerIP()
	if err != nil {
		return err
	}

	// Register TCP service
	tcpRegistration := new(api.AgentServiceRegistration)
	tcpRegistration.ID = fmt.Sprintf("file-server-tcp-%s", ip)
	tcpRegistration.Name = "file-server-tcp"
	tcpRegistration.Port = 8080
	tcpRegistration.Tags = []string{"tcp", "file", "server"}
	tcpRegistration.Address = ip

	// Register HTTP service
	httpRegistration := new(api.AgentServiceRegistration)
	httpRegistration.ID = fmt.Sprintf("file-server-http-%s", ip)
	httpRegistration.Name = "file-server-http"
	httpRegistration.Port = 8081
	httpRegistration.Tags = []string{"http", "file", "server"}
	httpRegistration.Address = ip

	// Add a health check for the HTTP service
	httpRegistration.Check = &api.AgentServiceCheck{
		HTTP:     fmt.Sprintf("http://%s:8081/health", ip),
		Interval: "10s",
		Timeout:  "1s",
	}

	// Register both services
	if err := client.Agent().ServiceRegister(tcpRegistration); err != nil {
		return err
	}
	if err := client.Agent().ServiceRegister(httpRegistration); err != nil {
		return err
	}

	return nil
}

func getContainerIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no IP address found")
}

// HealthCheckHandler handles health check requests
func (fs *FileServer) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
