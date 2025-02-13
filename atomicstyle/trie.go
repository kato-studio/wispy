package atomicstyle

// TrieNode represents a node in the trie.
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
	cssRule  string
}

// NewTrieNode creates a new trie node.
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
	}
}

// Trie is the main structure wrapping the trie root.
type Trie struct {
	root *TrieNode
}

// NewTrie creates a new trie.
func NewTrie() *Trie {
	return &Trie{root: NewTrieNode()}
}

// Insert adds a class and its CSS rule into the trie.
func (t *Trie) Insert(className, cssRule string) {
	node := t.root
	for _, char := range className {
		if _, exists := node.children[char]; !exists {
			node.children[char] = NewTrieNode()
		}
		node = node.children[char]
	}
	node.isEnd = true
	node.cssRule = cssRule
}

// Search looks up a class name and returns its CSS rule if found.
func (t *Trie) Search(className string) (string, bool) {
	node := t.root
	for _, char := range className {
		if child, exists := node.children[char]; exists {
			node = child
		} else {
			return "", false
		}
	}
	if node.isEnd {
		return node.cssRule, true
	}
	return "", false
}
