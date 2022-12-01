package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

func main() {
	open := flag.Bool("open", false, "On systems with an 'open' command, runs 'open $SERVER_URL'")
	dir := flag.String("path", ".", "The directory to serve via HTTP. Defaults to current directory.")
	flag.Parse()

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	addr := listener.Addr().String()
	fmt.Println("Running server:", addr)
	if *open {
		go openHost(addr)
	}
	http.Serve(listener, http.FileServer(http.Dir(*dir)))
}

func openHost(addr string) {
	url := url.URL{
		Scheme: "http",
		Host:   addr,
	}
	cmd := exec.Command("open", url.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
