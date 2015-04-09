package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
	}
	service := os.Args[1]

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	host := strings.Split(service, ":")[0]

	fmt.Println(host)

	header := "GET /index.htm HTTP/1.0 \r\n" + "Host:" + host + " \r\n\r\n"
	_, err = conn.Write([]byte(header))
	checkError(err)

	result, err := ioutil.ReadAll(conn)
	checkError(err)

	//	fmt.Println(string(result))

	content := strings.Split(string(result), "\r\n\r\n")[1]

	//	fmt.Println(string(content))

	err = ioutil.WriteFile("/Users/Eric/Desktop/Http_Server/downloaded/"+host+".html", []byte(content), os.ModePerm)
	checkError(err)

	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
