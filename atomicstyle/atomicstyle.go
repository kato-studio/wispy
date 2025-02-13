package atomicstyle

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// ExtractClasses parses HTML from the reader and extracts unique class names.
func ExtractClasses(r io.Reader) map[string]bool {
	classes := make(map[string]bool)
	tokenizer := html.NewTokenizer(r)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := tokenizer.Token()
		// Check for start or self-closing tags.
		if token.Type == html.StartTagToken || token.Type == html.SelfClosingTagToken {
			for _, attr := range token.Attr {
				if attr.Key == "class" {
					// Split classes by whitespace.
					for _, className := range strings.Fields(attr.Val) {
						classes[className] = true
					}
				}
			}
		}
	}

	return classes
}

// --- Prefix Mappings ---

var mediaQueryPrefixes = map[string]string{
	"sm":  "(min-width: 640px)",
	"md":  "(min-width: 768px)",
	"lg":  "(min-width: 1024px)",
	"xl":  "(min-width: 1280px)",
	"2xl": "(min-width: 1536px)",
}

var statePseudoPrefixes = map[string]string{
	"hover":  ":hover",
	"focus":  ":focus",
	"active": ":active",
	// Extend with additional state prefixes if needed.
}

// --- CSS Generation ---

// GenerateCSS accepts a set of class names and the trie, returning generated CSS.
func GenerateCSS(classes map[string]bool, trie *Trie) string {
	var buffer bytes.Buffer
	var rules []string

	for className := range classes {
		if rule, ok := generateRuleForClass(className, trie); ok {
			rules = append(rules, rule)
		}
	}

	sort.Strings(rules)
	for _, rule := range rules {
		buffer.WriteString(rule + "\n")
	}

	return buffer.String()
}

// GenerateDynamicCSS processes classes with dynamic values (currently supports w-[...] and h-[...])
func GenerateDynamicCSS(class string) (string, bool) {
	// Split the class on colon to separate any prefixes (e.g. media, state)
	parts := strings.Split(class, ":")
	base := parts[len(parts)-1]
	prefixes := parts[:len(parts)-1]

	// Example: handle dynamic width: "w-[500px]"
	if strings.HasPrefix(base, "w-[") && strings.HasSuffix(base, "]") {
		value := base[len("w-[") : len(base)-1]
		ruleBody := fmt.Sprintf("width: %s;", value)
		selector := buildSelector(class, prefixes)
		return wrapWithMediaQuery(selector, ruleBody, prefixes), true
	}
	// Example: handle dynamic height: "h-[500px]"
	if strings.HasPrefix(base, "h-[") && strings.HasSuffix(base, "]") {
		value := base[len("h-[") : len(base)-1]
		ruleBody := fmt.Sprintf("height: %s;", value)
		selector := buildSelector(class, prefixes)
		return wrapWithMediaQuery(selector, ruleBody, prefixes), true
	}
	// Add additional dynamic utilities as needed.
	return "", false
}

// generateRuleForClass checks the static trie first and then attempts dynamic generation.
func generateRuleForClass(class string, trie *Trie) (string, bool) {
	if rule, ok := trie.Search(class); ok {
		return "." + escapeClass(class) + " { " + rule + " }", true
	}
	// Fallback to dynamic generation (if supported)
	if dynRule, ok := GenerateDynamicCSS(class); ok {
		return dynRule, true
	}
	return "", false
}

// buildSelector constructs the CSS selector including state pseudo-classes.
func buildSelector(originalClass string, prefixes []string) string {
	// Use the full original class (with prefixes) for a unique selector and escape special characters.
	selector := "." + escapeClass(originalClass)
	// Append pseudo-classes for state prefixes (ignore media prefixes)
	for _, p := range prefixes {
		if pseudo, ok := statePseudoPrefixes[p]; ok {
			selector += pseudo
		}
	}
	return selector
}

// wrapWithMediaQuery wraps the rule in a media query if a media prefix is present.
func wrapWithMediaQuery(selector, ruleBody string, prefixes []string) string {
	for _, p := range prefixes {
		if mq, ok := mediaQueryPrefixes[p]; ok {
			innerRule := selector + " { " + ruleBody + " }"
			return fmt.Sprintf("@media %s { %s }", mq, innerRule)
		}
	}
	return selector + " { " + ruleBody + " }"
}

// escapeClass escapes special characters (such as colon and square brackets) in class names.
func escapeClass(class string) string {
	s := strings.ReplaceAll(class, ":", "\\:")
	s = strings.ReplaceAll(s, "[", "\\[")
	s = strings.ReplaceAll(s, "]", "\\]")
	return s
}

func Begin() string {
	trie := BuildExtendedTrie()
	input, _ := os.Open("./example.html")
	defer input.Close()

	fmt.Println("Execution Times")
	fmt.Println("------------------")
	extractTime := time.Now()
	// Extract unique class names from the HTML.
	classes := ExtractClasses(input)
	fmt.Println("Extract: ", time.Since(extractTime))

	generationTime := time.Now()
	// Generate CSS rules for the extracted classes.
	cssOutput := GenerateCSS(classes, trie)
	fmt.Println("Generate: ", time.Since(generationTime))
	fmt.Println("------------------")

	// Output the generated CSS.
	// fmt.Println(cssOutput)

	return cssOutput
}
