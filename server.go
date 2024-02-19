package main

import (
	"bufio"
	"flag"
	"fmt" // pass in command line
	"io"  // get file path
	"net"
	"os"
	"strings"
)

var (
	okResponse       = []byte("HTTP/1.1 200 OK\r\n\r\n")
	notFoundResponse = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
)

type Request struct {
	method  string
	path    string
	version string
	headers map[string]string
}

type Connection interface {
	io.Reader
	io.Writer
	Close() error
}

// TODO: add test

func main() {
	fmt.Println("Logs from program appear here")

	var dirFlag = flag.String("Directory", ".", "directory parsed from here ---")

	flag.Parse()

	if dirFlag == nil {
		fmt.Println("Error parsing command line argument.")
		os.Exit(1)
	}

	fmt.Println("dirFlat:", *dirFlag)

	listen, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221", err)
		os.Exit(1)
	}
	fmt.Print("Listening on port 4221")

	connection, err := listen.Accept()

	fmt.Println("Connection accepted")

	if err != nil {
		fmt.Println("Failed to accept connection", err.Error())
		os.Exit(1)
	}
	for {

		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {

	fmt.Println("Client connected")

	// Close the listener when the application closes.
	// use defer close to where the resource is created to ensure the resource is cleaned up as soon as it is no longer needed. This is to prevent later errors and easy to debug. Preventing us to forget to clean up the resource and memory leaks.
	defer connection.Close()
	// Read the incoming connection
	// TODO: bufio util explain
	connectionReader := bufio.NewReader(connection)

	// TODO: add ReadString explain

	requestInformation, err := connectionReader.ReadString('\n')

	if err != nil {
		fmt.Println("Failed to read request", err.Error())
		os.Exit(1)
	}

	// TODO: add strings.Fields explain
	requestInfoParts := strings.Fields(requestInformation)

	request := Request{
		// TODO: add request header, body, method, path, version explain
		method:  requestInfoParts[0],
		path:    requestInfoParts[1],
		version: requestInfoParts[2],
		headers: make(map[string]string),
	}

	for {
		headerLine, err := connectionReader.ReadString('\n')
		if err != nil || headerLine == "\r\n" {
			// TODO: this shouldn't break, handle with appropriate error response
			break
		}

		// TODO: splitN explain
		headerParts := strings.SplitN(headerLine, ": ", 2)
		if len(headerParts) == 2 {
			key := headerParts[0]
			value := headerParts[1]
			request.headers[key] = strings.TrimSpace(value)
		}
	}

	var response string

	// TODO: if else or switch case? is there any better way to handle this?
	if request.path == "/" {
		response = "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n<!DOCTYPE html><html><body><h1>Include content in response</h1></body></html>"

	} else if request.path == "/user-agent" {
		response = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nContent-Length: 11\r\n\r\ncurl/7.64.1"
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"

	}

	_, err = connection.Write([]byte(response))

	if err != nil {
		fmt.Println("Error writing data to connection:", err.Error())

		return
	}

}
