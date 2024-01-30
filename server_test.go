package main

import (
	"bytes"
	"fmt"
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

type TestCase struct {
	request  TestRequest
	expected TestResponse
}

type TestRequest struct {
	Method  string
	Path    string
	Host    string
	Headers map[string]string
}

type TestResponse struct {
	Status  string
	Headers map[string]string
	Body    string
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

func requestToString(req TestRequest) string {
	requestLine := fmt.Sprintf("%s %s HTTP/1.1\r\n", req.Method, req.Path)
	headers := "Host: " + req.Host + "\r\n"
	for name, value := range req.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", name, value)
	}
	return requestLine + headers + "\r\n"
}

func TestHandleConnection(t *testing.T) {
	tests := []TestCase{
		{
			request: TestRequest{
				Method:  "GET",
				Path:    "/",
				Host:    "localhost:4221",
				Headers: map[string]string{},
			},
			expected: TestResponse{
				Status:  "HTTP/1.1 200 OK\r\n",
				Headers: map[string]string{},
				Body:    "<!DOCTYPE html><html><body><h1>Include content in response</h1></body></html>",
			},
		},
		{
			request: TestRequest{
				Method:  "GET",
				Path:    "/user-agent",
				Host:    "localhost:4221",
				Headers: map[string]string{},
			},
			expected: TestResponse{
				Status:  "HTTP/1.1 200 OK\r\n",
				Headers: map[string]string{},
				Body:    "User-Agent: curl/7.64.1\r\n",
			},
		},
		{
			request: TestRequest{
				Method:  "GET",
				Path:    "/unknown-path",
				Host:    "localhost:4221",
				Headers: map[string]string{},
			},
			expected: TestResponse{
				Status:  "HTTP/1.1 404 Not Found\r\n",
				Headers: map[string]string{},
				Body:    "",
			},
		},
	}

	for _, test := range tests {

		t.Run(test.request.Path, func(t *testing.T) {
			input := requestToString(test.request)

			// Initialize a bytes.Buffer for the Data field in MockConnection
			dataBuffer := bytes.NewBufferString(input)

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
				t.Errorf("Connection not closed for input: %s", test.expected.Body)
			}
		})
	}

}
