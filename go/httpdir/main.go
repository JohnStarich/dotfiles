package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	dir := "."
	if len(os.Args) >= 2 {
		dir = os.Args[1]
	}
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting server:", listener.Addr().String())
	http.Serve(listener, http.FileServer(http.Dir(dir)))
}
