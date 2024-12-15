package main

import (
	"strings"
	"unicode"
)

type TokenType string

const (
	TokenTagOpen   TokenType = "TagOpen"
	TokenTagClose  TokenType = "TagClose"
	TokenAttribute TokenType = "Attribute"
	TokenContent   TokenType = "Content"
	TokenEOF       TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value string
}

func tokenize(input string) []Token {
	var tokens []Token
	var buffer strings.Builder

	for i := 0; i < len(input); i++ {
		ch := input[i]

		switch {
		case ch == '<': // Start of a tag
			if buffer.Len() > 0 {
				tokens = append(tokens, Token{Type: TokenContent, Value: buffer.String()})
				buffer.Reset()
			}
			tokens = append(tokens, Token{Type: TokenTagOpen, Value: "<"})

		case ch == '>': // End of a tag
			if buffer.Len() > 0 {
				tokens = append(tokens, Token{Type: TokenAttribute, Value: buffer.String()})
				buffer.Reset()
			}
			tokens = append(tokens, Token{Type: TokenTagClose, Value: ">"})

		case unicode.IsSpace(rune(ch)): // Ignore whitespace outside content
			if buffer.Len() > 0 {
				tokens = append(tokens, Token{Type: TokenAttribute, Value: buffer.String()})
				buffer.Reset()
			}

		default: // Collect content or attributes
			buffer.WriteByte(ch)
		}
	}

	if buffer.Len() > 0 {
		tokens = append(tokens, Token{Type: TokenContent, Value: buffer.String()})
	}

	tokens = append(tokens, Token{Type: TokenEOF, Value: ""})
	return tokens
}

type Node struct {
	Tag        string
	Attributes map[string]string
	Content    string
	Children   []Node
}

func TokenParse(tokens []Token) Node {
	var root Node
	stack := []Node{}

	for i := 0; i < len(tokens); i++ {
		switch tokens[i].Type {
		case TokenTagOpen:
			// Start a new tag
			tag := tokens[i+1].Value
			node := Node{Tag: tag, Attributes: make(map[string]string)}
			stack = append(stack, node)
			i++ // Skip over tag name

		case TokenAttribute:
			// Handle attributes
			attrParts := strings.SplitN(tokens[i].Value, "=", 2)
			if len(attrParts) == 2 {
				key := attrParts[0]
				value := strings.Trim(attrParts[1], `"`)
				stack[len(stack)-1].Attributes[key] = value
			}

		case TokenContent:
			// Add content to the current node
			stack[len(stack)-1].Content = tokens[i].Value

		case TokenTagClose:
			// Close the current tag
			node := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if len(stack) > 0 {
				// Add as a child to the parent node
				parent := &stack[len(stack)-1]
				parent.Children = append(parent.Children, node)
			} else {
				// It's the root node
				root = node
			}

		case TokenEOF:
			break
		}
	}

	return root
}
