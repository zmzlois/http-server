package main

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"
)

type MockConnection struct {
	io.Reader
	io.Writer
	closed        bool
	readDeadline  time.Time
	writeDeadline time.Time
	Data          *bytes.Buffer
}

func (mc *MockConnection) SetDeadline(t time.Time) error {
	mc.readDeadline = t
	return nil
}

func (mc *MockConnection) SetReadDeadline(t time.Time) error {
	mc.readDeadline = t
	return nil
}

func (mc *MockConnection) SetWriteDeadline(t time.Time) error {
	mc.writeDeadline = t
	return nil
}

func (mc *MockConnection) Read(b []byte) (n int, err error) {
	if time.Now().After(mc.readDeadline) {
		return 0, net.ErrClosed
	}
	return mc.Data.Read(b)
}

func (mc *MockConnection) Write(b []byte) (n int, err error) {
	if time.Now().After(mc.writeDeadline) {
		return 0, net.ErrClosed
	}
	return mc.Data.Write(b)
}

func (m *MockConnection) LocalAddr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 4221,
	}
}

func (mc *MockConnection) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8080,
	}
}

func (m *MockConnection) Close() error {
	m.closed = true
	return nil
}

func TestHandleConnection(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: "GET / HTTP/1.1\r\nHost: localhost:0.0.0.0:4221\r\n\r\n",

			expected: "HTTP/1.1 200 OK\r\n\r\n",
		},
		{
			input:    "GET /unknown-path HTTP/1.1\r\nHost: localhost:0.0.0.0:4221\r\n\r\n",
			expected: "HTTP/1.1 404 Not Found\r\n\r\n",
		},
	}

	for _, test := range tests {

		// Initialize a bytes.Buffer for the Data field in MockConnection
		dataBuffer := bytes.NewBufferString(test.input)

		conn := &MockConnection{
			Reader:        dataBuffer,
			Writer:        new(bytes.Buffer), // Add a Writer to the MockConnection
			Data:          dataBuffer,        // Set Data field to the buffer
			readDeadline:  time.Now().Add(1 * time.Second),
			writeDeadline: time.Now().Add(1 * time.Second),
		}

		// Set the read deadline to 1 second from now and prevent the connection from being closed too early
		conn.readDeadline = time.Now().Add(1 * time.Second)

		handleConnection(conn)

		if !conn.closed {
			t.Errorf("Connection not closed for input: %s", test.input)
		}
	}

}
