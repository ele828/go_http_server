package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

type Server struct {
	port string
	path string
}

var MimeMap = map[string]string{
	"txt":  "text/plain",
	"rtf":  "application/rtf",
	"pdf":  "application/pdf",
	"png":  "image/png",
	"gif":  "image/gif",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"html": "text/html",
}

func NewServer() *Server {
	return &Server{}
}

func (this *Server) SetPath(path string) {
	this.path = path
}

func (this *Server) Listen(port string) {
	this.port = port
	listener, err := net.Listen("tcp", "127.0.0.1:"+port)

	fmt.Println("Server running on http://127.0.0.1:" + this.port)

	defer listener.Close()

	if err != nil {
		fmt.Println("Error Listening:", err.Error())
	}

	for {
		conn, err := listener.Accept()
		checkError(err)

		go func(conn net.Conn) {
			uri, err := requestHandler(conn)
			if err == nil {
				this.routerHandler(uri, conn)
			}
		}(conn)
	}
}

// 服务器路由
func (this *Server) routerHandler(uri string, conn net.Conn) {

	if uri == "/" {
		uri = "/index.html"
	} else if uri == "/favicon.ico" {
		uri = "/index.html"
	}

	path := this.path + uri
	content, err := readStaticFile(path)
	mime := getMime(uri)

	// 打开失败或文件不存在，返回404页面
	if err != nil {
		content, err = readStaticFile(path + "index.html")
		mime = getMime("index.html")
		if err != nil {
			handle404(conn)
		}
	}

	header := "HTTP/1.0 200 OK\r\n" +
		"Server: Go_Server 1.0\r\n" +
		"Content-length: " + string(len(content)) + "\r\n" +
		"Connection: close\r\n" +
		"Content-type: " + mime + "\r\n\r\n"

	conn.Write([]byte(header))
	conn.Write([]byte(content))

	conn.Close()
}

// 请求处理
func requestHandler(conn net.Conn) (string, error) {

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		conn.Close()
		return "", err
	}
	req := string(buf)
	st := strings.Split(req, "GET ")
	token := strings.Split(st[1], "HTTP/1.1")[0]
	uri := strings.Trim(token, " ")
	//	fmt.Println(uri)
	return uri, nil
}

// 获取MIME类型
func getMime(file string) string {
	if strings.Contains(file, ".") {
		fileType := strings.Split(file, ".")[1]
		if mime, ok := MimeMap[fileType]; ok {
			return mime
		}
	}
	return "text/plain"
}

// 404页面
func handle404(conn net.Conn) {
	content := "404"
	header := "HTTP/1.0 404 NotFound\r\n" +
		"Server: Go_Server 1.0\r\n" +
		"Content-length: " + string(len(content)) + "\r\n" +
		"Connection: close\r\n" +
		"Content-type: " + "text/html" + "\r\n\r\n"

	conn.Write([]byte(header))
	conn.Write([]byte(content))

	conn.Close()
}

// 从服务器根目录读取文件
func readStaticFile(file string) ([]byte, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Error:%s", err.Error())
		os.Exit(1)
	}
}
