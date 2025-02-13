package atomicstyle

import (
	"bytes"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

//
// // use regex
// func RegexExtractClasses(htmlContent string
// func RegexExtractClasses(htmlContent string

// func RegexExtractClasses(htmlContent string) *dt.OrderedMap[string, struct{}] {
// 	classRegex := regexp.MustCompile(`class="([^"]+)"`)
// 	matches := classRegex.FindAllStringSubmatch(htmlContent, -1)
// 	classes := dt.NewOrderedMap[string, struct{}]()
// 	for _, match := range matches {
// 		for _, class_name := range strings.Split(match[1], " ") {
// 			classes.Set(class_name, struct{}{})
// 		}
// 	}
// 	return classes
// }

// func EscapeClassName(raw_class_name, state_string string) string {
// 	var escaped_class_name = raw_class_name
// 	var character_list = []string{":", ".", "]", "[", "/"}
// 	for _, character := range character_list {
// 		escaped_class_name = strings.ReplaceAll(escaped_class_name, character, "\\"+character)
// 	}
// 	return escaped_class_name + state_string
// }

// var CLASS_STATES = map[string]string{
// 	"hover":    "hover",
// 	"focus":    "focus",
// 	"active":   "active",
// 	"visited":  "visited",
// 	"disabled": "disabled",
// 	"first":    "first-child",
// 	"last":     "last-child",
// 	"odd":      "nth-child(odd)",
// 	"even":     "nth-child(even)",
// 	"after":    ":after",
// 	"before":   ":before",
// }

type Cache map[string]string
type Engine struct {
	Cache *Cache
}

// UtilityGenerator defines a function that, given a class name,
// returns a CSS rule (as a string) if the class is recognized.
type UtilityGenerator func(className string) (cssRule string, ok bool)

// registry holds all registered utility generators.
var (
	registryMu sync.RWMutex
	registry   []UtilityGenerator
)

// RegisterUtility adds a new UtilityGenerator to the registry.
func RegisterUtility(gen UtilityGenerator) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = append(registry, gen)
}

// GenerateCSSForClass iterates through the registry to see if any generator
// produces a CSS rule for the given class name.
func GenerateCSSForClass(className string) (string, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	for _, gen := range registry {
		if rule, ok := gen(className); ok {
			return rule, true
		}
	}
	return "", false
}

// ExtractClassesFromHTML parses the HTML from r and returns a set of unique class names.
func ExtractClassesFromHTML(r io.Reader) (map[string]struct{}, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	classes := make(map[string]struct{})
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "class" {
					for _, c := range strings.Fields(attr.Val) {
						classes[c] = struct{}{}
					}
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(doc)
	return classes, nil
}

// GenerateCSSFromClasses iterates over the unique class names and, using the registry,
// builds the complete CSS output.
func GenerateCSSFromClasses(classes map[string]struct{}) string {
	var buf bytes.Buffer
	// Sort class names for consistent output.
	keys := make([]string, 0, len(classes))
	for k := range classes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, className := range keys {
		if rule, ok := GenerateCSSForClass(className); ok {
			buf.WriteString(rule)
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

// init registers the default utility generators.
func init() {
	RegisterUtility(bgUtility)
	RegisterUtility(textUtility)
	RegisterUtility(mUtility)
	RegisterUtility(pUtility)
}

func Begin() string {
	input, _ := os.Open("./example.html")

	// Extract unique class names from the HTML.
	classes, err := ExtractClassesFromHTML(input)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	// Generate CSS rules for the extracted classes.
	cssOutput := GenerateCSSFromClasses(classes)

	// Output the generated CSS.
	// fmt.Println(cssOutput)

	return cssOutput
}
