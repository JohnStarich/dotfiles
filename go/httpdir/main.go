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
	basePath := flag.String("base", "/", "The base URL to use. All paths must use this prefix.")
	flag.Parse()

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	addr := listener.Addr().String()
	fmt.Println("Running server:", addr, *basePath)
	if *open {
		go openHost(addr, *basePath)
	}
	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir(*dir))
	mux.Handle(*basePath+"/", http.StripPrefix(*basePath, handler))
	http.Serve(listener, mux)
}

func openHost(addr, basePath string) {
	url := url.URL{
		Scheme: "http",
		Host:   addr,
		Path:   basePath,
	}
	cmd := exec.Command("open", url.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
