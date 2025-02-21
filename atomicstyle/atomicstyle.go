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

// Map for responsive breakpoints
var responsivePrefixes = map[string]string{
	"sm":  "@media (min-width: 640px)",
	"md":  "@media (min-width: 768px)",
	"lg":  "@media (min-width: 1024px)",
	"xl":  "@media (min-width: 1280px)",
	"2xl": "@media (min-width: 1536px)",
}

// For group variants, we’ll handle them specially in buildSelector.
var statePseudoPrefixes = map[string]string{
	"hover":             ":hover",
	"focus":             ":focus",
	"active":            ":active",
	"visited":           ":visited",
	"checked":           ":checked",
	"disabled":          ":disabled",
	"enabled":           ":enabled",
	"read-only":         ":read-only",
	"read-write":        ":read-write",
	"focus-within":      ":focus-within",
	"focus-visible":     ":focus-visible",
	"autofill":          ":autofill",
	"placeholder-shown": ":placeholder-shown",
	"default":           ":default",
	"first":             ":first-child",
	"last":              ":last-child",
	"only":              ":only-child",
	"odd":               ":nth-child(odd)",
	"even":              ":nth-child(even)",
	"first-of-type":     ":first-of-type",
	"last-of-type":      ":last-of-type",
	"only-of-type":      ":only-of-type",
	"empty":             ":empty",
	"open":              "[open]",
}

// --- Arbitrary Value Mapping ---
// Map a utility prefix (e.g. "w", "m", "px", etc.) to its corresponding CSS property.
var arbitraryMappings = map[string]string{
	// Spacing & Sizing
	"w":     "width",
	"h":     "height",
	"min-w": "min-width",
	"min-h": "min-height",
	"max-w": "max-width",
	"max-h": "max-height",
	"m":     "margin",
	"mt":    "margin-top",
	"mr":    "margin-right",
	"mb":    "margin-bottom",
	"ml":    "margin-left",
	"mx":    "margin-left; margin-right",
	"my":    "margin-top; margin-bottom",
	"p":     "padding",
	"pt":    "padding-top",
	"pr":    "padding-right",
	"pb":    "padding-bottom",
	"pl":    "padding-left",
	"px":    "padding-left; padding-right",
	"py":    "padding-top; padding-bottom",
	"gap":   "gap",
	"gap-x": "column-gap",
	"gap-y": "row-gap",

	// Colors
	"bg":     "background-color",
	"border": "border-color",
	"text":   "color",
	"fill":   "fill",
	"stroke": "stroke",

	// Borders & Outlines
	"rounded":  "border-radius",
	"border-w": "border-width",
	"border-t": "border-top-width",
	"border-r": "border-right-width",
	"border-b": "border-bottom-width",
	"border-l": "border-left-width",
	"ring":     "box-shadow",
	"outline":  "outline",

	// Positioning
	"z":       "z-index",
	"top":     "top",
	"right":   "right",
	"bottom":  "bottom",
	"left":    "left",
	"inset":   "inset",
	"inset-x": "left; right",
	"inset-y": "top; bottom",

	// Typography
	"leading":     "line-height",
	"tracking":    "letter-spacing",
	"font":        "font-family",
	"text-size":   "font-size",
	"text-weight": "font-weight",

	// Flexbox & Grid
	"flex":      "flex",
	"grow":      "flex-grow",
	"shrink":    "flex-shrink",
	"basis":     "flex-basis",
	"order":     "order",
	"grid-cols": "grid-template-columns",
	"grid-rows": "grid-template-rows",
	"col-span":  "grid-column",
	"row-span":  "grid-row",

	// Misc
	"opacity":    "opacity",
	"shadow":     "box-shadow",
	"blur":       "filter: blur",
	"brightness": "filter: brightness",
	"contrast":   "filter: contrast",
	"grayscale":  "filter: grayscale",
	"hue-rotate": "filter: hue-rotate",
	"invert":     "filter: invert",
	"saturate":   "filter: saturate",
	"sepia":      "filter: sepia",
}

// --- CSS Generation ---
// ExtractClasses parses HTML from the reader and extracts unique class names.
// ExtractClasses parses HTML from the reader and extracts class names in order.
func ExtractClasses(r io.Reader) []string {
	seen := make(map[string]bool) // Track unique class names
	var classList []string        // Preserve order

	tokenizer := html.NewTokenizer(r)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := tokenizer.Token()
		if token.Type == html.StartTagToken || token.Type == html.SelfClosingTagToken {
			for _, attr := range token.Attr {
				if attr.Key == "class" {
					for _, className := range strings.Fields(attr.Val) {
						if !seen[className] { // Ensure uniqueness while keeping order
							seen[className] = true
							classList = append(classList, className)
						}
					}
				}
			}
		}
	}
	return classList
}

// Updated function to generate CSS dynamically
// func GenerateDynamicCSS(class string) (string, bool) {
// 	parts := strings.Split(class, ":")
// 	base := parts[len(parts)-1]
// 	prefixes := parts[:len(parts)-1]

// 	if i := strings.Index(base, "-["); i != -1 && strings.HasSuffix(base, "]") {
// 		propPrefix := base[:i]
// 		arbitraryValue := base[i+2 : len(base)-1]
// 		if rule, ok := generateArbitraryRule(propPrefix, arbitraryValue); ok {
// 			selector := buildSelector(class, prefixes)
// 			return selector + " { " + rule + " }", false
// 		}
// 	}
// 	return "", false
// }

// generateArbitraryRule converts a utility prefix and its arbitrary value to a CSS rule.
// func generateArbitraryRule(prefix, value string) (string, bool) {
// 	if cssProp, ok := arbitraryMappings[prefix]; ok {
// 		// For compound properties (like "mx" → "margin-left; margin-right")
// 		if strings.Contains(cssProp, ";") {
// 			parts := strings.Split(cssProp, ";")
// 			var rules []string
// 			for _, p := range parts {
// 				p = strings.TrimSpace(p)
// 				if p != "" {
// 					rules = append(rules, fmt.Sprintf("%s: %s;", p, value))
// 				}
// 			}
// 			return strings.Join(rules, " "), true
// 		}
// 		return fmt.Sprintf("%s: %s;", cssProp, value), true
// 	}
// 	return "", false
// }

// --- Selector & Media Wrapping Helpers ---

// buildSelector constructs the CSS selector using the full original class (with proper escaping)
// and then applies state pseudo-classes and group variants based on the provided prefixes.
func BuildSelector(originalClass string, prefixes []string) (selector string, mediaQuery string) {
	selector = "." + escapeClass(originalClass)
	var pseudoClasses []string
	var groupSelector string

	for _, p := range prefixes {
		if strings.HasPrefix(p, "group-") {
			// Handle group-based variants
			switch p {
			case "group-hover":
				groupSelector = ".group:hover "
			case "group-focus":
				groupSelector = ".group:focus "
			case "group-active":
				groupSelector = ".group:active "
			case "group-aria-expanded":
				groupSelector = ".group[aria-expanded='true'] "
			case "group-aria-selected":
				groupSelector = ".group[aria-selected='true'] "
			}
		} else if strings.HasPrefix(p, "peer-") {
			// Handle peer-based variants
			switch p {
			case "peer-hover":
				groupSelector = ".peer:hover ~ "
			case "peer-focus":
				groupSelector = ".peer:focus ~ "
			case "peer-checked":
				groupSelector = ".peer:checked ~ "
			case "peer-disabled":
				groupSelector = ".peer:disabled ~ "
			}
		} else if strings.HasPrefix(p, "aria-") {
			// Handle ARIA attributes
			pseudoClasses = append(pseudoClasses, fmt.Sprintf("[aria-%s='true']", strings.TrimPrefix(p, "aria-")))
		} else if strings.HasPrefix(p, "data-") {
			// Handle data attributes
			pseudoClasses = append(pseudoClasses, fmt.Sprintf("[%s]", p))
		} else if mq, exists := responsivePrefixes[p]; exists {
			// Handle responsive prefixes
			mediaQuery = mq
		} else if pseudo, ok := statePseudoPrefixes[p]; ok {
			pseudoClasses = append(pseudoClasses, pseudo)
		} else if strings.HasPrefix(p, "not-") {
			// Handle "not-" pseudo-classes
			notState := strings.TrimPrefix(p, "not-")
			if pseudo, exists := statePseudoPrefixes[notState]; exists {
				pseudoClasses = append(pseudoClasses, fmt.Sprintf(":not(%s)", pseudo))
			}
		}
	}

	if groupSelector != "" {
		selector = groupSelector + selector
	}
	if len(pseudoClasses) > 0 {
		selector += strings.Join(pseudoClasses, "")
	}
	return selector, mediaQuery
}

func generateRuleForClass(class string, trie *Trie) (rule string, mediaQuery string, ok bool) {
	parts := strings.Split(class, ":")
	base := parts[len(parts)-1]      // get last item
	prefixes := parts[:len(parts)-1] // all except last item

	// Try a static lookup on the base utility.
	if ruleBody, ok := trie.Search(base); ok {
		selector, mediaQuery := BuildSelector(class, prefixes)
		return selector + " { " + ruleBody + " }", mediaQuery, true
	} else {
		fmt.Println("[-]", base)
	}

	// Fallback: if not found statically, attempt dynamic generation on the full class.
	// if dynRule, ok := GenerateDynamicCSS(class); ok {
	// 	return dynRule, true
	// }
	return "", "", false
}

// GenerateCSS accepts a set of class names and the trie, returning generated CSS.
func GenerateCSS(classes []string, trie *Trie) string {
	var buffer bytes.Buffer
	var defaultRules []string
	mediaQueryRules := make(map[string][]string)
	var mediaQueryList []string // Maintain insertion order

	for _, className := range classes {
		if rule, mediaQuery, ok := generateRuleForClass(className, trie); ok {
			if mediaQuery == "" {
				// Rules without media queries go into the default group
				defaultRules = append(defaultRules, rule)
			} else {
				// If it's the first time seeing this media query, track order
				if _, exists := mediaQueryRules[mediaQuery]; !exists {
					mediaQueryList = append(mediaQueryList, mediaQuery)
				}
				mediaQueryRules[mediaQuery] = append(mediaQueryRules[mediaQuery], rule)
			}
		}
	}

	// Output default rules (no media query)
	for _, rule := range defaultRules {
		buffer.WriteString(rule + "\n")
	}

	// Sort media queries by Tailwind priority
	sort.SliceStable(mediaQueryList, func(i, j int) bool {
		return mediaQueryPriority(mediaQueryList[i]) < mediaQueryPriority(mediaQueryList[j])
	})

	// Output rules grouped by ordered media queries
	for _, mq := range mediaQueryList {
		buffer.WriteString(mq + " {\n")
		for _, rule := range mediaQueryRules[mq] {
			buffer.WriteString("  " + rule + "\n")
		}
		buffer.WriteString("}\n")
	}

	return buffer.String()
}

// Define priority for Tailwind-like media queries
func mediaQueryPriority(mq string) int {
	priority := map[string]int{
		"@media (min-width: 640px)":  1, // sm
		"@media (min-width: 768px)":  2, // md
		"@media (min-width: 1024px)": 3, // lg
		"@media (min-width: 1280px)": 4, // xl
		"@media (min-width: 1536px)": 5, // 2xl
	}
	if p, exists := priority[mq]; exists {
		return p
	}
	return 99 // Default lowest priority for unknown media queries
}

// escapeClass escapes special characters (such as colon and square brackets) in class names.
func escapeClass(class string) string {
	s := strings.ReplaceAll(class, ":", "\\:")
	s = strings.ReplaceAll(s, "[", "\\[")
	s = strings.ReplaceAll(s, "]", "\\]")
	s = strings.ReplaceAll(s, ".", "\\.")
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "/", "\\/")
	return s
}

func BeginTest() string {
	trie := BuildFullTrie()
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
