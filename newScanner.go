package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kato-studio/wispy/engine"
)

// Raw HTML input
var rawHTMLSmall = `
boop
<TitleThing title="I Am A Title!">
	lorem ipsum dolor sit amet consectetur adipiscing elit
</TitleThing>
<section class="section section--resources section--faq" id="">
	      	<Card class="container  ">
		inner card content stuff here with nested header above!
		TEST
		<Card class="container  ">
			<TitleThing title="I Am A Title!"> prime lorem ipsum dolor sit amet consectetur adipiscing elit</TitleThing>
			inner card content stuff here with nested header above!
		</Card>
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

func newScannerV2(w http.ResponseWriter, req *http.Request) {
	var result strings.Builder
	var scanOne_time = time.Now()
	//
	scanner := engine.NewScanner(strings.NewReader(string(rawHTMLSmall)))
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
