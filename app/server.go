package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// func sendResponse(res string, conn net.Conn) {
// 	conn.Write([]byte(res + "\n"))
// }

// func handleRequest(conn net.Conn) {
// 	for {
// 		log.Println("Handling Request")
// 		buffer := make([]byte, 1024)

// 		length, err := conn.Read(buffer)
// 		if err != nil {
// 			log.Panicln(err)
// 		}

// 		str := string(buffer[:length])

// 		fmt.Println(conn.RemoteAddr().String())
// 		fmt.Printf("Received command %d\t:%s\n", length, str)

// 		switch str {
// 		case "PING\r\n":
// 			sendResponse("PONG", conn)
// 		case "PUSH\r\n":
// 			sendResponse("GOT PUSH", conn)
// 		case "QUIT\r\n":
// 			sendResponse("Goodbye", conn)
// 			conn.Close()
// 		default:
// 			conn.Write([]byte(fmt.Sprintf("UNKNOWN_COMMAND: %s\n", str)))
// 		}
// 	}
// }

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	var conn, error = l.Accept()
	if error != nil {
		fmt.Println("Error accepting connection: ", error.Error())
		os.Exit(1)
	} else {
		buffer := make([]byte, 1024)
		length, err := conn.Read(buffer)
		if err != nil {
			log.Panicln(err)
		}

		str := string(buffer[:length])

		result := strings.Split(str, " ")

		// fmt.Println(result)

		path := result[1]

		// fmt.Println(conn.RemoteAddr().String())
		// fmt.Printf("Received command %d\t:%s\n", length, str)

		switch {
		case path == "/":
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		case strings.HasPrefix(path, "/echo"):
			pathQuery := strings.Split(path, "/echo/")
			resultStr := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(pathQuery[1])) + "\r\n\n" + pathQuery[1] + "\r\n"
			conn.Write([]byte(resultStr))
		case strings.HasPrefix(path, "/user-agent"):
			userAgentInfo := strings.Split(str, "User-Agent: ")
			userAgentInfo = strings.Split(userAgentInfo[1], "\n")
			fmt.Println(userAgentInfo[0])
			resultStr := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(userAgentInfo[0])) + "\r\n\n" + userAgentInfo[0] + "\r\n"
			conn.Write([]byte(resultStr))
		default:
			conn.Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"))
		}
	}

	conn.Close()
}
