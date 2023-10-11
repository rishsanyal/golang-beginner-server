package main

import (
	"fmt"
	"log"
	"net"
	"os"
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

		path := result[1]

		// fmt.Println("Here")
		// fmt.Println(path)

		fmt.Println(conn.RemoteAddr().String())
		fmt.Printf("Received command %d\t:%s\n", length, str)

		if path == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else {
			conn.Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"))
		}
	}

	conn.Close()
}
