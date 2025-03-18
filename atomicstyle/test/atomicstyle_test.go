package atomicstyle_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kato-studio/wispy/atomicstyle"
)

func BeginTest() string {
	trie := atomicstyle.BuildFullTrie()
	input, _ := os.Open("./example.html")
	defer input.Close()

	fmt.Println("Execution Times")
	fmt.Println("------------------")
	extractTime := time.Now()
	// Extract unique class names from the HTML.
	classes := atomicstyle.ExtractClasses(input)
	fmt.Println("Extract: ", time.Since(extractTime))
	generationTime := time.Now()
	// Generate CSS rules for the extracted classes.
	cssOutput := atomicstyle.GenerateCSS(classes, trie)
	fmt.Println("Generate: ", time.Since(generationTime))
	fmt.Println("------------------")

	return cssOutput
}

func TestBegin(t *testing.T) {
	result := BeginTest()
	log.Printf("Generated CSS: \n%s", result)

	if result == "" {
		t.Errorf("Expected generated CSS, but got empty string")
	}
}
