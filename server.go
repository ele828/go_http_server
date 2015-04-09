package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var first = true

func main() {

	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	for {
		conn, err := l.Accept()
		checkError(err)
		go func(conn net.Conn) {
			uri := RequestHandler(conn)
			routerHandler(uri, conn)
			return
		}(conn)
	}
}

func RequestHandler(conn net.Conn) string {
	//	fmt.Println(first)
	if first == true {
		first = false
		time.Sleep(time.Second * 10)
	}

	//	fmt.Println(conn.RemoteAddr())

	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	req := string(buf)
	st := strings.Split(req, "GET ")
	token := strings.Split(st[1], "HTTP/1.1")[0]
	uri := strings.Trim(token, " ")
	//	fmt.Println(uri)
	return uri
}

func routerHandler(uri string, conn net.Conn) {

	if uri == "/" {
		uri = "index.html"
	} else if uri == "/favicon.ico" {
		uri = "index.html"
	}

	//	content, err := readStaticFile("/Users/Eric/Desktop/Http_Server/htdocs/" + uri)
	//	checkError(err)
	content := "testing"

	header := "HTTP/1.0 200 OK\r\n" +
		"Server: Go_Server 1.0\r\n" +
		"Content-length: " + string(len(content)) + "\r\n" +
		"Connection: close" +
		"Content-type: text/html\r\n\r\n"

	conn.Write([]byte(header))
	conn.Write([]byte(content))

	conn.Close()
	return
}

func readStaticFile(file string) ([]byte, error) {
	f, err := ioutil.ReadFile(file)
	checkError(err)
	return f, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Error:%s", err.Error())
		os.Exit(1)
	}
}

func test(conn net.Conn) {
	fmt.Println(first)
	if first == true {
		first = false
		time.Sleep(time.Second * 10)
	}
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	header := "HTTP/1.0 200 OK\r\n" +
		"Server: Go_Server 1.0\r\n" +
		"Content-length: " + string(len("abc")) + "\r\n" +
		"Content-type: text/html\r\n\r\n"

	conn.Write([]byte(header))

	conn.Write([]byte("test"))

	conn.Close()

	return
}
