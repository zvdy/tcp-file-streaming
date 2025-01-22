package http

import (
	"crypto/rand"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/zvdy/tcp-file-streaming/internal/tcp"
)

type HTTPServer struct {
	Port string
}

func (hs *HTTPServer) Start() {
	http.HandleFunc("/health", hs.HealthCheckHandler)
	http.HandleFunc("/stream", hs.StreamHandler)
	log.Println("HTTP server started on port", hs.Port)
	if err := http.ListenAndServe(hs.Port, nil); err != nil {
		log.Fatal("HTTP server error:", err)
	}
}

func (hs *HTTPServer) StreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Stream request headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("%s: %s\n", name, value)
		}
	}

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

	file := make([]byte, size)
	_, err = io.ReadFull(rand.Reader, file)
	if err != nil {
		http.Error(w, "Failed to generate file data", http.StatusInternalServerError)
		return
	}

	err = tcp.SendToTCPServer(file)
	if err != nil {
		http.Error(w, "Failed to send data to TCP server", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data processed successfully"))
}

func (hs *HTTPServer) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
