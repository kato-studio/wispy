package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/kato-studio/wispy/engine"
	"github.com/kato-studio/wispy/engine/xops"
)

var site = "/abc-example.com"

// Testing Page Rendering with Atomic CSS
func abc(w http.ResponseWriter, req *http.Request) {
	var json_start = time.Now()
	fmt.Println("JSON in", time.Since(json_start))
	var ctx_start = time.Now()

	var ctx = engine.NewCtx(engine.TemplateCtx{
		Page: engine.TemplatePageCtx{
			Title: "ABC Example",
			Head: "<link rel=\"stylesheet\" href=\"/css/style.css\">" +
				"",
			Meta:    "<meta name=\"description\" content=\"ABC Example\">",
			Styles:  "",
			Scripts: "",
		},
		Data: map[string]any{
			"abc": "example",
			"Site": map[string]any{
				"Title": "ABC Example",
			},
		},
	})
	fmt.Println("Context in", time.Since(ctx_start))

	var ren_start = time.Now()
	var r = engine.InitEngine(ctx)
	var page_html, err = r.RenderPage(site, "")
	fmt.Println("Rendered in", time.Since(ren_start))
	fmt.Println("Total in", time.Since(json_start))

	if err != nil {
		fmt.Println(err)
	}
	io.WriteString(w, page_html)
}

// Raw HTML input
var rawHTMLSmall = `

`

func init() {
	// rawHTMLSmall = strings.Repeat(rawHTMLSmall, 3)
	file, err := os.ReadFile("./test.html")
	if err != nil {
		fmt.Print(err)
	}
	rawHTMLSmall = string(file)
}

func NewScannerV2(w http.ResponseWriter, req *http.Request) {
	var scanOne_time = time.Now()
	var ctx = engine.NewCtx(engine.TemplateCtx{
		Page: engine.TemplatePageCtx{
			Title:   "ABC Example",
			Head:    "<link rel=\"stylesheet\" href=\"/css/style.css\">",
			Meta:    "<meta name=\"description\" content=\"ABC Example\">",
			Styles:  "",
			Scripts: "",
		},
	})

	var r = engine.InitEngine(ctx)
	// Add the operation function to the OperationFuncMap
	r.SetOperationFunc("each", xops.EachOperation)
	fmt.Print("\n\n\n\n\n\n\n\n\n\n")
	fmt.Println("restarting..")
	fmt.Print("\n\n\n\n\n\n\n\n\n\n")
	result, err := r.Html(rawHTMLSmall)

	if err != nil {
		fmt.Println("\033[31m", err, "\033[0m")
	}

	scanOne_Duration := time.Since(scanOne_time)

	fmt.Println("\n\nScan Duration: ", scanOne_Duration)
	w.Header().Add("content-type", "text/html")
	io.WriteString(w, result)
}

func main() {
	const port = ":8090"
	http.HandleFunc(site, abc)
	http.HandleFunc("/scanner", NewScannerV2)
	fmt.Print("Listening on port:", port)
	http.ListenAndServe(port, nil)
}
