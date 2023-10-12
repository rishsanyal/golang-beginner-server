package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func handleFileResponse(conn net.Conn, directoryLocn string, requestStr string) string {
	fmt.Println(requestStr)
	pathStr := strings.TrimSpace(strings.Split(requestStr, "/files/")[1])

	pathStr = strings.Split(pathStr, " ")[0]
	// fmt.Println(pathStr)
	resultStr := ""

	if _, err := os.Stat(directoryLocn + pathStr); err == nil {

		fileStream, err := ioutil.ReadFile(directoryLocn + pathStr)
		if err != nil {
			log.Panicln(err)
			fmt.Println(err.Error())
			return err.Error()
		}

		// fmt.Println(fileStream)
		resultStr = "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + strconv.Itoa(len(fileStream)) + "\r\n\n" + string(fileStream) + "\r\n"
	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		// fmt.Println(fileStream)
		resultStr = "HTTP/1.1 404 NOT FOUND\r\n\r\n"
	} else {
		log.Fatalf("In the Else issue")
	}

	return resultStr
}

func handlePostFileResponse(conn net.Conn, directoryLocn string, requestStr string) string {
	pathStr := strings.TrimSpace(strings.Split(requestStr, "/files/")[1])

	pathStr = strings.Split(pathStr, " ")[0]
	// fmt.Println(pathStr)
	fileBodyArr := strings.Split(requestStr, "\n")
	fileBody := fileBodyArr[len(fileBodyArr)-2]

	// resultStr := ""

	fmt.Println(fileBody)

	d1 := []byte(fileBody + "\r\n")
	err := os.WriteFile((directoryLocn + pathStr), d1, 0644)

	if err != nil {
		log.Panic(err.Error())
		// resultStr
	}

	resultStr := "HTTP/1.1 201 Created\r\n\r\n"

	// if _, err := os.Stat(directoryLocn + pathStr); err == nil {

	// 	fileStream, err := ioutil.ReadFile(directoryLocn + pathStr)
	// 	if err != nil {
	// 		log.Panicln(err)
	// 		fmt.Println(err.Error())
	// 		return err.Error()
	// 	}

	// 	// fmt.Println(fileStream)
	// 	resultStr = "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + strconv.Itoa(len(fileStream)) + "\r\n\n" + string(fileStream) + "\r\n"
	// } else if errors.Is(err, os.ErrNotExist) {
	// 	// path/to/whatever does *not* exist
	// 	// fmt.Println(fileStream)
	// 	resultStr = "HTTP/1.1 404 NOT FOUND\r\n\r\n"
	// } else {
	// 	log.Fatalf("In the Else issue")
	// }

	return resultStr
}

func handleConnection(conn net.Conn, directoryLocn string) {

	buffer := make([]byte, 1024)
	length, err := conn.Read(buffer)
	if err != nil {
		log.Panicln(err)
	}

	str := string(buffer[:length])

	result := strings.Split(str, " ")

	// fmt.Println(result[0])

	requestType := result[0]
	path := result[1]

	switch {
	case requestType == "GET" && path == "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	case requestType == "GET" && strings.HasPrefix(path, "/echo"):
		pathQuery := strings.Split(path, "/echo/")
		resultStr := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(pathQuery[1])) + "\r\n\n" + pathQuery[1] + "\r\n"
		conn.Write([]byte(resultStr))
	case requestType == "GET" && strings.HasPrefix(path, "/user-agent"):
		userAgentInfo := strings.Split(str, "User-Agent: ")
		userAgentInfo = strings.Split(userAgentInfo[1], "\n")
		tempUserAgentInfo := strings.TrimSpace(userAgentInfo[0])
		resultStr := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(tempUserAgentInfo)) + "\r\n\n" + tempUserAgentInfo + "\n"
		conn.Write([]byte(resultStr))
	case requestType == "GET" && strings.HasPrefix(path, "/files"):
		resultStr := handleFileResponse(conn, directoryLocn, str)
		conn.Write([]byte(resultStr))
	case requestType == "POST" && strings.HasPrefix(path, "/files"):
		resultStr := handlePostFileResponse(conn, directoryLocn, str)
		conn.Write([]byte(resultStr))
	default:
		conn.Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"))
	}

	conn.Close()
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	argsSlice := os.Args[1:]
	var directoryLocn string = ""

	if len(argsSlice) > 0 && argsSlice[0] == "--directory" {
		directoryLocn = argsSlice[1]
	} else {
		directoryLocn = ""
	}

	fmt.Println(directoryLocn)

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {

		var conn, error = l.Accept()
		if error != nil {
			fmt.Println("Error accepting connection: ", error.Error())
			os.Exit(1)
		} else {
			go handleConnection(conn, directoryLocn)
		}

	}

}
