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

	defer connection.Close()

	connectionReader := bufio.NewReader(connection)

	requestInformation, err := connectionReader.ReadString('\n')

	if err != nil {
		fmt.Println("Failed to read request", err.Error())
		os.Exit(1)
	}

	requestInfoParts := strings.Fields(requestInformation)

	request := Request{}

	request.method = requestInfoParts[0]
	request.path = requestInfoParts[1]
	request.version = requestInfoParts[2]
	request.headers = make(map[string]string)

	for {
		headerLine, err := connectionReader.ReadString('\n')
		if err != nil || headerLine == "\r\n" {
			break
		}
		headerParts := strings.SplitN(headerLine, ": ", 2)
		if len(headerParts) == 2 {
			key := headerParts[0]
			value := headerParts[1]
			request.headers[key] = strings.TrimSpace(value)
		}
	}

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