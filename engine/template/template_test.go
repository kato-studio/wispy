package template

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

// renderFromFileTest demonstrates rendering a template from file.
func TestRenderFromFileTest(t *testing.T) {
	// Open the grammar file (for validation/documentation purposes).
	file, err := os.Open("template.ebnf")
	if err != nil {
		log.Fatalf("Error opening grammar file: %v", err)
	}
	defer file.Close()

	// Create a new template engine.
	engine := NewTemplateEngine()

	// Set up the rendering context using NewRenderCtx (which initializes Internal automatically).
	ctx := NewRenderCtx(map[string]interface{}{
		"Domain":      "example.com",
		"title":       "Welcome",
		"showContent": "true", // any non-empty string except "false" is truthy
		"content":     "   This is some sample content.   ",
		"items":       "apple, banana, cherry",
		"condition":   "true",
	})

	// Render the sample template from file with timing logs.
	fmt.Println("\nRendered Template:")
	testFilePath := filepath.Join("mock", "test.html")
	output, err := RenderFile(testFilePath, engine, ctx, true)
	if err != nil {
		log.Fatalf("Error rendering template from file: %v", err)
	}

	// Write the output to an HTML file.
	buf := bytes.NewBufferString(output)
	if err := os.WriteFile("output.html", buf.Bytes(), 0777); err != nil {
		log.Printf("Error writing output to file: %v", err)
	}

	fmt.Println("------------------")
	fmt.Println()

	// Large test: attempt to read from large.html if available.
	largeFilePath := filepath.Join("mock", "large.html")
	if _, err := os.Stat(largeFilePath); err == nil {
		log.Printf("%sRunning large test using %s%s", colorBlue, largeFilePath, colorReset)
		fmt.Println("\nRendered Large Template:")
		output, err := RenderFile(largeFilePath, engine, ctx, false)
		if err != nil {
			log.Printf("%sError rendering large template: %v%s", colorRed, err, colorReset)
		} else {
			_ = output
			// Optionally print or process the large output.
		}
	} else {
		log.Printf("%sLarge test file %s not found, skipping large test.%s", colorYellow, largeFilePath, colorReset)
	}
}
