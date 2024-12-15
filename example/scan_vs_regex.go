package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/tdewolff/parse/html"
)

/*
- -------------------------------
-  TESTING CODE
- -------------------------------
*/

var componentRegistry = map[string]string{}

func registerComponent(name, template string) {
	componentRegistry[name] = template
}

func extractAttribute(attributes, name string) string {
	re := regexp.MustCompile(name + `="(.*?)"`)
	match := re.FindStringSubmatch(attributes)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func renderWithRegex(rawHTML string) string {
	re := regexp.MustCompile(`<([a-zA-Z]+)(.*?)>(.*?)</{1}>`)
	return re.ReplaceAllStringFunc(rawHTML, func(match string) string {
		matches := re.FindStringSubmatch(match)
		tag := matches[1]
		attributes := matches[2]
		content := matches[3]

		if template, exists := componentRegistry[tag]; exists {
			rendered := template
			// Replace placeholders like {{title}} and {{content}}
			rendered = strings.ReplaceAll(rendered, "{{title}}", extractAttribute(attributes, "title"))
			rendered = strings.ReplaceAll(rendered, "{{content}}", content)
			return rendered
		}
		return match
	})
}

func extractAttributes(line string) map[string]string {
	attrs := make(map[string]string)
	re := regexp.MustCompile(`(\w+)="([^"]+)"`)
	matches := re.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		attrs[match[1]] = match[2]
	}
	return attrs
}

func extractContent(line, tag string) string {
	re := regexp.MustCompile(fmt.Sprintf(`<%s[^>]*>(.*?)</%s>`, tag, tag))
	match := re.FindStringSubmatch(line)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func renderWithScanner(rawHTML string) string {
	scanner := bufio.NewScanner(strings.NewReader(rawHTML))
	var result strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line contains a component
		if strings.Contains(line, "<") && strings.Contains(line, ">") {
			startTagIndex := strings.Index(line, "<")
			endTagIndex := strings.Index(line, ">")
			if startTagIndex >= 0 && endTagIndex > startTagIndex {
				tag := line[startTagIndex+1 : endTagIndex]
				if template, exists := componentRegistry[tag]; exists {
					attributes := extractAttributes(line)
					content := extractContent(line, tag)
					rendered := strings.ReplaceAll(template, "{{title}}", attributes["title"])
					rendered = strings.ReplaceAll(rendered, "{{content}}", content)
					result.WriteString(rendered)
				} else {
					result.WriteString(line)
				}
			} else {
				result.WriteString(line)
			}
		} else {
			result.WriteString(line)
		}
		result.WriteString("\n")
	}

	return result.String()
}

func AverageTime(xs []int64) int64 {
	total := int64(0)
	for _, x := range xs {
		total += x
	}
	return total / int64(len(xs))
}
func measureTime(name string, fn func() string) time.Duration {
	start := time.Now()
	fn()
	elapsed := time.Since(start)
	var x = elapsed
	// return fmt.Sprintf("[%s] Execution Time: %s\n\n", name, elapsed) + fn() + "\n\n"
	// fmt.Println(fmt.Sprintf("[%s] Execution Time: %d\n\n", name, x))
	return x
}

func scannerTest(w http.ResponseWriter, req *http.Request) {

	// Register reusable components
	registerComponent("TitleThing", `<div class="title"><h1>{{title}}</h1><small class="sub">{{content}}</small></div>`)
	registerComponent("Card", `<div class="card"><h2>{{title}}</h2><div class="card-body">{{content}}</div></div>`)

	// Raw HTML input
	rawHTML := `
		<TitleThing title="I Am A Title!">
			lorem ipsum dolor sit amet consectetur adipiscing elit
		</TitleThing>
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

		<section class="section section--resources section--faq" id="">
    <div class="container  ">
      <h2 class="section__title  " tabindex="0"></h2>

      <p class="section__desc" tabindex="0"></p>

      <div class="section__content" tabindex="0">

        <div class="row">
          <div class="col-lg-5 mob-m-t-0x mobile-center" style="display: block;">
            <h3>
              Full Stack Cloud Compute
            </h3>
            <div>
              Build a composable cloud according to your business needs.
            </div><br>
            <div class="flex">
              <a class="btn btn--sm btn--outline" style="z-index: 1; color: var(--brand-primary);"
                href="/products/cloud-compute/">
                <span class="btn__text">Learn more</span>
                <span class="btn-hover-effect"></span>
              </a>
            </div>
          </div>
          <div class="order-1 order-lg-2 col-lg-7">
            <div class="box box--faq" style="min-height: auto;" data-faq-container="">
              <div class="box__faq" style="min-width: 270px;">
                <div class="box__faq-body">
                  <div class="list-group">
                    <div class="list-group__item list-group--collapse">
                      <div class="list-group__top top  is-active " data-faq-question="">
                        <div class="top__title" data-faq-question-title="">
                          Cloud Compute
                        </div>
                      </div>
                      <div class="list-group__collapse show " data-faq-answer="">
                        <div class="list-group__content" data-faq-answer-content="">
                          <p>
                            Build a composable cloud according to your business needs.
                          </p>
                          <div class="text-center flex">
                            <a class="btn btn--sm btn--primary" style="z-index: 1;" href="/products/cloud-compute/">
                              <span class="btn__text">Learn more</span>
                              <span class="btn-hover-effect"></span></a>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="list-group__item list-group--collapse">
                      <div class="list-group__top top " data-faq-question="">
                        <div class="top__title" data-faq-question-title="">
                          Optimized Cloud Compute
                        </div>
                      </div>
                      <div class="list-group__collapse" data-faq-answer="">
                        <div class="list-group__content" data-faq-answer-content="">
                          <p>
                            Powered by the latest-generation AMD and Intel CPUs, spin up virtual machines for general
                            purpose or optimized configurations in under 60 seconds.
                          </p>
                          <div class="text-center flex">
                            <a class="btn btn--sm btn--primary" style="z-index: 1;"
                              href="/products/optimized-cloud-compute/">
                              <span class="btn__text">Learn more</span>
                              <span class="btn-hover-effect"></span></a>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="list-group__item list-group--collapse">
                      <div class="list-group__top top " data-faq-question="">
                        <div class="top__title" data-faq-question-title="">
                          Bare Metal
                        </div>
                      </div>
                      <div class="list-group__collapse" data-faq-answer="">
                        <div class="list-group__content" data-faq-answer-content="">
                          <p>
                            Stay in full control of your environment with high performance single-tenant dedicated
                            servers, accelerated by NVIDIA GPUs and high-performance CPUs from AMD and Intel.
                          </p>
                          <div class="text-center flex">
                            <a class="btn btn--sm btn--primary" style="z-index: 1;" href="/products/bare-metal/">
                              <span class="btn__text">Learn more</span>
                              <span class="btn-hover-effect"></span></a>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="list-group__item list-group--collapse">
                      <div class="list-group__top top " data-faq-question="">
                        <div class="top__title" data-faq-question-title="">
                          Kubernetes
                        </div>
                      </div>
                      <div class="list-group__collapse" data-faq-answer="">
                        <div class="list-group__content" data-faq-answer-content="">
                          <p>
                            Deploy and scale containerized apps with a fully managed service. Vultr Kubernetes Engine
                            ushers in a better way to optimize container orchestration.
                          </p>
                          <div class="text-center flex">
                            <a class="btn btn--sm btn--primary" style="z-index: 1;" href="/kubernetes/">
                              <span class="btn__text">Learn more</span>
                              <span class="btn-hover-effect"></span></a>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="list-group__item list-group--collapse">
                      <div class="list-group__top top " data-faq-question="">
                        <div class="top__title" data-faq-question-title="">
                          Managed Databases
                        </div>
                      </div>
                      <div class="list-group__collapse" data-faq-answer="">
                        <div class="list-group__content" data-faq-answer-content="">
                          <p>
                            Secure, highly available, and easily scalable, Vultr Managed Databases for MySQL,
                            PostgreSQL, Apache Kafka®, and Redis®*-compatible Vultr Managed Databases for Caching 'just
                            work' right out of the box.
                          </p>
                          <div class="text-center flex">
                            <a class="btn btn--sm btn--primary" style="z-index: 1;" href="/products/managed-databases/">
                              <span class="btn__text">Learn more</span>
                              <span class="btn-hover-effect"></span></a>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="list-group__item list-group--collapse">
                      <div class="list-group__top top " data-faq-question="">
                        <div class="top__title" data-faq-question-title="">
                          Storage
                        </div>
                      </div>
                      <div class="list-group__collapse" data-faq-answer="">
                        <div class="list-group__content" data-faq-answer-content="">
                          <p>
                            With Vultr, block storage and object storage solutions are flexible, scalable, and
                            expandable without sacrificing performance or security.
                          </p>
                          <div class="text-center flex">
                            <a class="btn btn--sm btn--primary" style="z-index: 1;" href="/products/block-storage/">
                              <span class="btn__text">Learn more</span>
                              <span class="btn-hover-effect"></span></a>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="box__answer">
                <h4 class="h6" data-faq-answer-title-view="">Cloud Compute</h4>
                <div data-faq-answer-view="">
                  <p>Build a composable cloud according to your business needs.</p>
                  <div class="text-center flex">
                    <a class="btn btn--sm btn--primary" style="z-index: 1;" href="/products/cloud-compute/">
                      <span class="btn__text">Learn more</span>
                      <span class="btn-hover-effect"></span></a>
                  </div>
                </div>
              </div>
            </div>
            <div class="redis-disclaimer">*Redis® is a registered trademark of Redis Ltd. Any rights therein are
              reserved to Redis Ltd. Any use by Vultr is for referential purposes only and does not indicate any
              sponsorship, endorsement or affiliation between Redis® and Vultr.</div>
          </div>
        </div>

      </div>





    </div>
  </section>
		<Card title="Dynamic Card Title 3">
			This is the card body content.
		</Card>
		<TitleThing title="I Am A Title! [Four]">
			lorem ipsum dolor sit amet consectetur adipiscing elit
		</TitleThing>
		<Card title="Dynamic Card Title 4">
			This is the card body content.
		</Card>
		<TitleThing title="I Am A Title! [Five]">
			lorem ipsum dolor sit amet consectetur adipiscing elit
		</TitleThing>
		<Card title="Dynamic Card Title 5">
			This is the card body content.
		</Card>

	`

	// Measure both parsing methods
	// CSV formate
	res := ""
	res += "Regex, Lexer, HtmlToken\n"

	for i := 0; i < 100; i++ {
		regex := measureTime("Regex Parsing", func() string { return renderWithRegex(rawHTML) })
		lexer := measureTime("Scanner Parsing", func() string { return renderWithScanner(rawHTML) })
		html_token := measureTime("Tokenize HTML", func() string {
			// Tokenize HTML from stdin.
			l := html.NewLexer(strings.NewReader(rawHTML))
			output := ""
			for {
				tt, data := l.Next()
				switch tt {
				case html.ErrorToken:
					if l.Err() != io.EOF {
						fmt.Println("Error on line", l.Err())
					}
					return ""
				case html.StartTagToken:
					for {
						ttAttr, _ := l.Next()
						if ttAttr != html.AttributeToken {
							break
						}

						// key := dataAttr
						// val := l.AttrVal()
						// fmt.Println("[", string(key), "=", string(val), "]")
						// output += fmt.Sprint("[", string(key), "=", string(val), "]")
					}
				case html.EndTagToken:
					output += fmt.Sprintf("#END#%s", data)
				case html.TextToken:
					output += string(data)
				// case html.CommentToken:
				// output += ""
				case html.DoctypeToken:
					// output += fmt.Sprintf("<!DOCTYPE %s>", data)
					output += string(data)
				default:
					output += fmt.Sprintf("**!!**%s**!!**\n", data)
				}
			}
		})
		res += fmt.Sprintf("%s, %s, %s\n", regex, lexer, html_token)
	}
	w.Write([]byte(res))
	// res += fmt.Sprintf("Regex: %dμs, Lexer: %dμs, HtmlToken: %dμs\n", AverageTime(regex), AverageTime(lexer), AverageTime(html_token))

	// res += measureTimeAndRender("Regex Parsing", func() string { return renderWithRegex(rawHTML) })
	// res += measureTimeAndRender("Scanner Parsing", func() string { return renderWithScanner(rawHTML) })
	// res += measureTimeAndRender("Tokenize HTML", func() string {
	// 	// Tokenize HTML from stdin.
	// 	l := html.NewLexer(parse.NewInput(strings.NewReader(rawHTML)))
	// 	output := ""
	// 	for {
	// 		tt, data := l.Next()
	// 		switch tt {
	// 		case html.ErrorToken:
	// 			if l.Err() != io.EOF {
	// 				fmt.Println("Error on line", l.Err())
	// 			}
	// 			return output
	// 		case html.StartTagToken:
	// 			for {
	// 				ttAttr, dataAttr := l.Next()
	// 				if ttAttr != html.AttributeToken {
	// 					break
	// 				}

	// 				key := dataAttr
	// 				val := l.AttrVal()
	// 				fmt.Println("[", string(key), "=", string(val), "]")
	// 				output += fmt.Sprint("[", string(key), "=", string(val), "]")
	// 			}
	// 		case html.EndTagToken:
	// 			output += fmt.Sprintf("#END#%s", data)
	// 		case html.TextToken:
	// 			output += string(data)
	// 		// case html.CommentToken:
	// 		// output += ""
	// 		case html.DoctypeToken:
	// 			// output += fmt.Sprintf("<!DOCTYPE %s>", data)
	// 			output += string(data)
	// 		default:
	// 			output += fmt.Sprintf("**!!**%s**!!**\n", data)
	// 		}
	// 	}
	// })
	// res += measureTimeAndRender("Lexer Parsing", func() string {
	// 	lexer := NewLexer(rawHTML)
	// 	fmt.Println("Tokens:")
	// 	output := ""
	// 	for {
	// 		token := lexer.NextLexToken()
	// 		output += fmt.Sprintf("%+v\n", token)
	// 		if token.Type == LexTokenEOF {
	// 			break
	// 		}
	// 	}
	// 	return output
	// })
	// res += measureTimeAndRender("Tokenize Parsing", func() string {
	// 	// Tokenize
	// 	tokens := tokenize(rawHTML)
	// 	return fmt.Sprintf("%+v\n", tokens)
	// })

	w.Write([]byte(res))
}
