package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Request struct {
	method  string
	path    string
	version string
	headers map[string]string
}

// TODO: add test

func main() {
	fmt.Println("Logs from program appear here")

	listen, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Filaed to bind to port 4221", err)
		os.Exit(1)
	}

	connection, err := listen.Accept()

	if err != nil {
		fmt.Println("Failed to accept connection", err.Error())
		os.Exit(1)
	}

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

	request := Request{}

	// TODO: add request header, body, method, path, version explain
	request.method = requestInfoParts[0]
	request.path = requestInfoParts[1]
	request.version = requestInfoParts[2]
	request.headers = make(map[string]string)

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

	// TODO: if else or switch case? is there any better way to handle this?
	if request.path == "/" {
		response := "HTTP/1.1 200 OK\r\n\r\n"
		_, err = connection.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing data to connection:", err.Error())
			os.Exit(1)
		}
	} else {
		response := "HTTP/1.1 404 Not Found\r\n\r\n"
		_, err = connection.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing data to connection:", err.Error())
			os.Exit(1)
		}
	}

}
