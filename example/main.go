package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kato-studio/wispy/engine"
)

func main() {
	const port = ":8090"
	http.HandleFunc("/html", TestHtml)
	http.HandleFunc("/tree", TestTreeGen)
	fmt.Print("Listening on port:", port)
	http.ListenAndServe(port, nil)
}

func TestHtml(res http.ResponseWriter, req *http.Request) {
	fmt.Println("\n\n\nRunning TestHtml :) ")
	var Full_time = time.Now()
	var output []byte = make([]byte, 0)
	var file []byte = make([]byte, 0)
	var err error
	var r engine.Render

	var attFuncMap = make(map[string]engine.AttributeFunc)
	attFuncMap["class"] = func(name, value string) (bool, string, []error) {
		fmt.Println(value)
		return false, "", []error{}
	}

	var opFuncMap = make(map[string]engine.OperationFunc)
	var pageCtx = engine.NewCtx(engine.TemplateCtx{
		Page: engine.TemplatePageCtx{
			Title: "ABC Example",
			Head:  "<link rel=\"stylesheet\" href=\"/css/style.css\">",
			Meta:  "<meta name=\"description\" content=\"ABC Example\">",
			Css:   "",
			Js:    "",
		},
		Data: map[string]any{},
	})
	r = engine.InitEngine(pageCtx, attFuncMap, opFuncMap)
	//
	fmt.Println("Init Duration: ", time.Since(Full_time))
	//
	var jsOutput strings.Builder
	for _, v := range r.Js {
		jsOutput.Write(v)
	}
	pageCtx.Page.Js += jsOutput.String()
	//
	var cssOutput strings.Builder
	for _, v := range r.Css {
		cssOutput.Write(v)
	}
	pageCtx.Page.Css += cssOutput.String()
	//
	var OsRead_Duration = time.Now()
	file, err = os.ReadFile("./test.html")
	fmt.Println("os-Read Duration: ", time.Since(OsRead_Duration))
	// Failed to html file
	if err != nil {
		// res.Write([]byte(err.Error()))
		res.Write([]byte("Failed to read file \n\n"))
		res.Write([]byte(err.Error()))
		return
	}

	Html_time := time.Now()
	localData := map[string]any{
		"foo":    "bar",
		"whitty": "whimsy",
	}
	output, errs := r.Render(file, localData, nil, nil, "page")
	//
	fmt.Println("Render Duration: ", time.Since(Html_time))
	//
	for _, er := range errs {
		fmt.Println(er)
	}
	//
	fmt.Println("Scan Duration: ", time.Since(Full_time))
	//
	// Headers
	res.Header().Add("content-type", "text/html")
	// res.Header().Add("content-type", "text/html")
	// res.Header().Add("content-type", "application/json")
	res.Write(output)
}

func TestTreeGen(res http.ResponseWriter, req *http.Request) {
	fmt.Println("\n\n\nRunning TestTreeGen()")
	var Full_time = time.Now()
	var output []byte = make([]byte, 0)
	var file []byte = make([]byte, 0)
	var namedImportPath = make(map[string]string)
	var err error
	var r engine.Render

	//
	fmt.Println("Init Duration: ", time.Since(Full_time))
	var OsRead_Duration = time.Now()
	file, err = os.ReadFile("./test.html")
	fmt.Println("os-Read Duration: ", time.Since(OsRead_Duration))
	// Failed to html file
	if err != nil {
		fmt.Println("[Error] ", err.Error())
		res.Write([]byte("Failed to read file"))
		res.Write([]byte(err.Error()))
		return
	}

	Html_time := time.Now()
	nodeTree, errs := r.BuildNodeTree(file, "page", namedImportPath)

	//
	fmt.Println("Tree Gen Duration: ", time.Since(Html_time))
	// log errors
	for _, er := range errs {
		fmt.Println(er.Error())
	}

	fmt.Println("Scan Duration: (no marshal)", time.Since(Full_time))
	output, err = json.Marshal(nodeTree)
	fmt.Println("Scan Duration: (after marshal)", time.Since(Full_time))
	if err != nil {
		fmt.Printf(err.Error())
		res.Write([]byte("500"))
	}

	//
	// Headers
	res.Header().Add("content-type", "application/json")
	res.Write(output)
}
