package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/kato-studio/wispy/atomic"
	"github.com/kato-studio/wispy/engine"
)

var site = "/abc-example.com"

func abc(w http.ResponseWriter, req *http.Request) {
	var json_start = time.Now()
	fmt.Println("JSON in", time.Since(json_start))
	var ctx_start = time.Now()
	var ctx = engine.NewTemplateCtx(engine.TemplateCtx{
		Page: engine.TemplatePageCtx{
			Title:      "ABC Example",
			Layout:     "./shared/layouts/_default" + engine.LAYOUT_FILE,
			InsertHead: "",
			Meta: []string{
				"name=\"description\" content=\"ABC Example\"",
				"name=\"keywords\" content=\"ABC, Example\"",
			},
			Css: "",
			Js:  "",
		},
	})
	fmt.Println("Context in", time.Since(ctx_start))
	//
	var ren_start = time.Now()
	var page_html, page_err = engine.RenderFile("./sites"+site+"/pages/"+engine.PAGE_FILE, ctx)
	if page_err != nil {
		fmt.Println(page_err)
	}
	fmt.Println("Rendered in", time.Since(ren_start))
	//
	var css_start = time.Now()
	compiled_css := atomic.Compile(page_html.String(), atomic.WispyStaticStyles, atomic.WispyColors)
	ctx.Page.Css += "\n---------- Atomic Styles ----------\n" + compiled_css + "\n-------------------------\n"
	fmt.Println("Atomic Css in", time.Since(css_start))
	//
	var page_with_layout, layout_err = engine.RenderPage(page_html, ctx)
	if layout_err != nil {
		fmt.Println(layout_err)
	}
	//
	w.Write(page_with_layout.Bytes())
	fmt.Println("Total in", time.Since(json_start))
	fmt.Println("------------------------------------------------")
}

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

func main() {
	http.HandleFunc(site, abc)
	http.HandleFunc("/self-closing", selfClosing)
	http.ListenAndServe(":8090", nil)
}
