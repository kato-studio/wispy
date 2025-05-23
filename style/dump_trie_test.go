package style

import (
	"fmt"
	"os"
	"testing"
)

func TestWriteTrieNodeToFile(t *testing.T) {
	// Get the large trie structure
	atomicTrie := BuildFullTrie()

	// Open a file for writing
	file, dumpErr := os.Create("trie_dump.txt")
	if dumpErr != nil {
		fmt.Println("Error creating file:", dumpErr)
		return
	}
	defer file.Close()

	// Write the trie structure to the file
	atomicTrie.Dump(file)
	fmt.Println("Trie structure written to trie_dump.txt")

	// Write the CSS rules to a file
	cssRrr := atomicTrie.WriteCSSToFile("output.css", BuildSelector)
	if cssRrr != nil {
		fmt.Println("Error writing CSS to file:", cssRrr)
		return
	}
	fmt.Println("CSS written to output.css")

}
