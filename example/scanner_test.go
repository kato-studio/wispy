package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/tdewolff/parse/html"
)

var __test_componentRegistry = map[string]string{}

func __test_registerComponent(name, template string) {
	__test_componentRegistry[name] = template
}

func __test_renderWithRegex(rawHTML string) string {
	re := regexp.MustCompile(`<([a-zA-Z]+)(.*?)>(.*?)</{1}>`)
	return re.ReplaceAllStringFunc(rawHTML, func(match string) string {
		matches := re.FindStringSubmatch(match)
		tag := matches[1]
		attributes := matches[2]
		content := matches[3]

		if template, exists := __test_componentRegistry[tag]; exists {
			rendered := template
			// Replace placeholders like {{title}} and {{content}}
			rendered = strings.ReplaceAll(rendered, "{{title}}", extractAttribute(attributes, "title"))
			rendered = strings.ReplaceAll(rendered, "{{content}}", content)
			return rendered
		}
		return match
	})
}

func __test_renderWithScanner(rawHTML string) string {
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

func __test_renderTokenizer(rawHTML string) string {
	// Tokenize HTML from stdin.
	l := html.NewLexer(strings.NewReader(rawHTML))
	output := ""
	for {
		tt, data := l.Next()
		switch tt {
		case html.ErrorToken:
			if l.Err() != io.EOF {
				// fmt.Println("Error on line", l.Err())
			}
			return output
		case html.StartTagToken:
			for {
				ttAttr, dataAttr := l.Next()
				if ttAttr != html.AttributeToken {
					break
				}

				key := dataAttr
				val := l.AttrVal()
				// fmt.Println("[", string(key), "=", string(val), "]")
				output += fmt.Sprint("[", string(key), "=", string(val), "]")
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
}

// func __test_measureTime(name string, fn func() string) int64 {
// 	start := time.Now()
// 	fn()
// 	return time.Since(start).Microseconds()
// }

// func __test_AverageTime(times []int64) int64 {
// 	var total int64
// 	for _, t := range times {
// 		total += t
// 	}
// 	return total / int64(len(times))
// }

var rawHTML = `
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

func Benchmark_RegexParsing(b *testing.B) {

	// Register reusable components
	__test_registerComponent("TitleThing", `<div class="title"><h1>{{title}}</h1><small class="sub">{{content}}</small></div>`)
	__test_registerComponent("Card", `<div class="card"><h2>{{title}}</h2><div class="card-body">{{content}}</div></div>`)

	for i := 0; i < b.N; i++ {
		__test_renderWithRegex(rawHTML)
	}
}

func Benchmark_ScannerParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		__test_renderWithScanner(rawHTML)
	}
}

func Benchmark_TokenizeHTML(b *testing.B) {
	for i := 0; i < b.N; i++ {
		__test_renderTokenizer(rawHTML)
	}
}
