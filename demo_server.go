package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kato-studio/wispy/engine"
	"github.com/tidwall/gjson"
)

var site = "/abc-example.com"

// Testing Page Rendering with Atomic CSS
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

// Raw HTML input
var rawHTMLSmall = `
boop
{{palceholder}}
<TitleThing title="I Am A Title!">
	lorem ipsum dolor sit amet consectetur adipiscing elit
</TitleThing>
<section class="section section--resources section--faq" id="">
{{palceholder}}
	      	<Card class="container  ">
		inner card content stuff here with nested header above!
		{{palceholder}}
		
		<Card class="container  ">
			<TitleThing title="I Am A Title!"> prime lorem ipsum dolor sit amet consectetur adipiscing elit</TitleThing>
			inner card content stuff here with nested header above!
		</Card>
		{{Eq("I'm an operation!")}}
		TEST END
	</Card>


	               <TitleThing title="I Am A Title!"> prime lorem ipsum dolor sit amet consectetur adipiscing elit</TitleThing>
	<Card title="Dynamic Card Title 0">
		This is the card body content.
	            </Card>
	<TitleThing title="I Am A Title! [Two]">
		lorem ipsum dolor sit amet consectetur adipiscing elit
	</TitleThing>
	<Card title="Dynamic Card Title 2">
		This is the card body content.
	</Card>
	<TitleThing title="I Am A Title! [Three]">
		lorem ipsum dolor sit amet consectetur adipiscing elit
	</TitleThing>
</section>
<section class="section section--resources section--faq" id="">
	{{Eq("I'm an operation!")}}

	{{For .list as link }}

	{{/For}}
	<div class="container  ">
		<h2 class="section__title  " tabindex="0"></h2>
		<p class="section__desc" tabindex="0">boop 2</p>
		<div class="section__content" tabindex="0"></div>
	</div>
</section>
`

func init() {
	// rawHTMLSmall = strings.Repeat(rawHTMLSmall, 3)
}

func NewScannerV2(w http.ResponseWriter, req *http.Request) {
	var result strings.Builder
	var scanOne_time = time.Now()
	//
	// scanner := engine.NewScanner(strings.NewReader(string(rawHTMLSmall)))
	s := engine.
		result.WriteString(`
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2.0.6/css/pico.classless.min.css">
	<style>
		code {
			display: block;
			width: fit-content;
			margin: 10px 0px;
		}
	</style>
	`)
	for scanner.Scan() {
		data := scanner.Bytes()
		// length := len(bytes)
		// if length > 3 {
		//
		// fir_byte := bytes[0]
		// sec_byte := bytes[1]
		//
		// if fir_byte == '<' && isCapitalByte(sec_byte) {
		// 	// result.WriteString("<Component:")
		// 	continue
		// }
		// }

		// isCapitalByte
		result.WriteString(`<code>` + html.EscapeString(string(data)) + `</code>`)
	}
	scanOne_Duration := time.Since(scanOne_time)
	fmt.Println("Scan Duration: ", scanOne_Duration)
	w.Header().Add("content-type", "text/html")
	io.WriteString(w, result.String())
}

func main() {
	const port = ":8090"
	http.HandleFunc(site, abc)
	http.HandleFunc("/scanner", NewScannerV2)
	http.ListenAndServe(port, nil)
	// err := http.ListenAndServe(port, nil)
	// if err != nil {
	// 	fmt.Println("Error", fmt.Sprint(err))
	// } else {
	// 	fmt.Println("Server listening on port:", port)
	// }

}
