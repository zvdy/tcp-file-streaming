package server

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"
)

type mockConn struct {
	net.Conn
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
}

func (m *mockConn) Read(b []byte) (int, error) {
	return m.readBuffer.Read(b)
}

func (m *mockConn) Write(b []byte) (int, error) {
	return m.writeBuffer.Write(b)
}

func (m *mockConn) Close() error {
	return nil
}

func TestHandleConnection(t *testing.T) {
	tests := []struct {
		name           string
		inputData      []byte
		expectedOutput string
	}{
		{
			name:           "Successful data transfer",
			inputData:      append(encodeSize(5), []byte("hello")...),
			expectedOutput: "done",
		},
		{
			name:           "Connection closed by client",
			inputData:      []byte{},
			expectedOutput: "",
		},
		{
			name:           "Error during data transfer",
			inputData:      encodeSize(5),
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readBuffer := bytes.NewBuffer(tt.inputData)
			writeBuffer := new(bytes.Buffer)
			conn := &mockConn{readBuffer: readBuffer, writeBuffer: writeBuffer}

			fs := &FileServer{}
			fs.handleConnection(conn)

			if tt.expectedOutput == "" {
				if writeBuffer.Len() != 0 {
					t.Errorf("expected no output, got %s", writeBuffer.String())
				}
			} else {
				if writeBuffer.String() != tt.expectedOutput {
					t.Errorf("expected %s, got %s", tt.expectedOutput, writeBuffer.String())
				}
			}
		})
	}
}

func encodeSize(size int64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, size)
	return buf.Bytes()
}
