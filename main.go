package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/skratchdot/open-golang/open"
)

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("reading stdin error:", err)
		os.Exit(1)
	}

	http.HandleFunc("/stdin", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", http.DetectContentType(input))
		resp.Write(input)
	})

	time.AfterFunc(1*time.Second, func() {
		fmt.Println("open stdin...")
		open.Run("http://localhost/stdin")
	})

	if err = http.ListenAndServe("", nil); err != nil {
		fmt.Println("web-viewer error:", err)
		os.Exit(2)
	}
}
