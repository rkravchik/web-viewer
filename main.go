package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"

	"github.com/skratchdot/open-golang/open"
)

// exit codes:
// 0 - successful
// 1 - os specific errors (stdin, exec etc)
// 2 - net errors (socket, server, http)
// 3 - program has been interrupted or killed by signal

func main() {
	code := app()
	os.Exit(code)
}

var (
	endpointBase  = "http://localhost"
	endpointStdin = "/stdin"
)

func app() int {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("reading stdin error:", err)
		return 1
	}

	closeChan := make(chan os.Signal, 1)
	signal.Notify(closeChan, os.Interrupt, os.Kill)

	completeChan := make(chan struct{})
	http.HandleFunc(endpointStdin, func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", http.DetectContentType(input))
		resp.Write(input)
		close(completeChan)
	})

	// don't use ListenAndServe to "catch" browser connection early
	srv := new(http.Server)
	defer srv.Close()
	// listen on local ip address only
	ln, err := net.Listen("tcp", "localhost:http")
	if err != nil {
		fmt.Println("open listen socket error:", err)
		return 2
	}

	fmt.Println("open stdin...")
	if err = open.Run(path.Join(endpointBase, endpointStdin)); err != nil {
		fmt.Println("opening stdin error:", err)
		return 1
	}

	errChan := make(chan error)
	go func() { errChan <- srv.Serve(ln) }()

	select {
	case <-completeChan:
	case err = <-errChan:
		fmt.Println("web-viewer error:", err)
		return 2
	case <-closeChan:
		fmt.Println("signal to close")
		return 3
	}

	return 0
}
