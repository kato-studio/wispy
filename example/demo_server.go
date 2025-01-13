package main

import (
	"fmt"
	"net/http"

	"github.com/kato-studio/wispy/engine"
)

var site = "/abc-example.com"

func main() {
	const port = ":8090"
	http.HandleFunc("/html", engine.TestHtml)
	http.HandleFunc("/tree", engine.TestTreeGen)
	fmt.Print("Listening on port:", port)
	http.ListenAndServe(port, nil)
}
