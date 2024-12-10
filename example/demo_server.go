package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kato-studio/wispy/engine"
	"github.com/tidwall/gjson"
)

var site = "/abc-example.com"

func abc(w http.ResponseWriter, req *http.Request) {
	var json_start = time.Now()
	json_parsed := gjson.Parse(`{
			"abc": "example",
			"Site": {
				"Title": "ABC Example",
			}
		}`)
	fmt.Println("JSON in", time.Since(json_start))
	var ctx_start = time.Now()
	var ctx = engine.NewTemplateCtx(engine.TemplateCtx{
		Page: engine.TemplatePageCtx{
			Title:   "ABC Example",
			Head:    "<link rel=\"stylesheet\" href=\"/css/style.css\">",
			Meta:    "<meta name=\"description\" content=\"ABC Example\">",
			Styles:  "",
			Scripts: "",
		},
		Json: json_parsed,
	})
	fmt.Println("Context in", time.Since(ctx_start))

	var ren_start = time.Now()
	var page_html, err = engine.RenderPage("./sites"+site+"/pages/"+engine.PAGE_FILE, ctx)
	fmt.Println("Rendered in", time.Since(ren_start))
	fmt.Println("Total in", time.Since(json_start))

	if err != nil {
		fmt.Println(err)
	}
	w.Write(page_html.Bytes())
}

func main() {
	http.HandleFunc(site, abc)
	http.ListenAndServe(":8090", nil)
}
