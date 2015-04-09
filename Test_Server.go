package main

import (
	"./server"
)

func main() {
	svr := server.NewServer()
	svr.SetPath("./htdocs")
	svr.Listen("8888")
}
