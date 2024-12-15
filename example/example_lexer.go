package main

import (
	"unicode"
)

// LexTokenType represents the type of a token.
type LexTokenType string

const (
	LexTokenTagOpen    LexTokenType = "TagOpen"
	LexTokenTagClose   LexTokenType = "TagClose"
	LexTokenIdentifier LexTokenType = "Identifier"
	LexTokenAttribute  LexTokenType = "Attribute"
	LexTokenString     LexTokenType = "String"
	LexTokenContent    LexTokenType = "Content"
	LexTokenEOF        LexTokenType = "EOF"
)

// LexToken represents a lexical token.
type LexToken struct {
	Type  LexTokenType
	Value string
}

// Lexer represents the lexical analyzer.
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

// NewLexer initializes a new lexer with the given input.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar advances the lexer to the next character.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // Indicates EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextLexToken produces the next token from the input.
func (l *Lexer) NextLexToken() LexToken {
	l.skipWhitespace()

	switch l.ch {
	case '<':
		l.readChar()
		if l.ch == '/' {
			l.readChar()
			return LexToken{Type: LexTokenTagClose, Value: "</"}
		}
		return LexToken{Type: LexTokenTagOpen, Value: "<"}
	case '>':
		l.readChar()
		return LexToken{Type: LexTokenTagClose, Value: ">"}
	case '"':
		return LexToken{Type: LexTokenString, Value: l.readString()}
	case 0:
		return LexToken{Type: LexTokenEOF, Value: ""}
	default:
		if isLetter(l.ch) {
			identifier := l.readIdentifier()
			if l.peekChar() == '=' {
				l.readChar() // Skip '='
				value := l.readString()
				return LexToken{Type: LexTokenAttribute, Value: identifier + `="` + value + `"`}
			}
			return LexToken{Type: LexTokenIdentifier, Value: identifier}
		}
		content := l.readContent()
		return LexToken{Type: LexTokenContent, Value: content}
	}
}

// Helper Functions

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readString() string {
	l.readChar() // Skip opening quote
	start := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	value := l.input[start:l.position]
	l.readChar() // Skip closing quote
	return value
}

func (l *Lexer) readContent() string {
	start := l.position
	for l.ch != '<' && l.ch != 0 {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '-'
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}
