package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
)

func selfClosing(w http.ResponseWriter, req *http.Request) {
	// load bightml.html file
	//
	raw_html, _ := os.ReadFile("bightml.html")
	replace_start := time.Now()
	var re = regexp.MustCompile(`<(\w*?) (.*?)\s*\/>`)
	result := re.ReplaceAllString(string(raw_html), `<$1 $2></$1>`)
	//
	fmt.Println("Replace in", time.Since(replace_start))
	w.Write([]byte(result))
}
