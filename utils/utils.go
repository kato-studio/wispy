package utils

import (
	"fmt"
	"regexp"
	"strings"
)

/* https://patorjk.com/software/taag/#p=testall&f=ANSI%20Regular&t=Logger
██       ██████   ██████   ██████  ███████ ██████
██      ██    ██ ██       ██       ██      ██   ██
██      ██    ██ ██   ███ ██   ███ █████   ██████
██      ██    ██ ██    ██ ██    ██ ██      ██   ██
███████  ██████   ██████   ██████  ███████ ██   ██
*/

type Logger struct {
	Print       func(any)    // Print logs
	Fatal       func(any)    // Fatal logs
	Error       func(string) // Error logs
	Warn        func(string) // Warn logs
	Info        func(string) // Info logs
	ServerPrint func(string) // ServerPrint logs
}

var logger Logger = Logger{
	Print: func(value interface{}) {
		fmt.Println(value)
	},
	Fatal: func(value interface{}) {
		fmt.Println(value)
	},
	Error: func(value string) {
		fmt.Println(value)
	},
	Warn: func(value string) {
		fmt.Println(value)
	},
	Info: func(value string) {
		fmt.Println(value)
	},
	ServerPrint: func(str string) {
		fmt.Println("----------------")
		fmt.Println(str)
		fmt.Println("----------------")
	},
}

func GetLogger() Logger {
	return logger
}

/* ----------------------
██    ██ ███    ██ ██  ██████  ██    ██ ███████ ███████ ███████ ████████
██    ██ ████   ██ ██ ██    ██ ██    ██ ██      ██      ██         ██
██    ██ ██ ██  ██ ██ ██    ██ ██    ██ █████   ███████ █████      ██
██    ██ ██  ██ ██ ██ ██ ▄▄ ██ ██    ██ ██           ██ ██         ██
 ██████  ██   ████ ██  ██████   ██████  ███████ ███████ ███████    ██
                          ▀▀
---------------------- */

// UniqueSet represents a set of unique strings
type UniqueSet map[string]struct{}

// NewUniqueSet creates and returns a new UniqueSet
func NewUniqueSet(initial ...string) UniqueSet {
	this_set := make(UniqueSet)
	for _, str := range initial {
		this_set.Add(str)
	}
	return this_set
}

// Add adds a string to the set
func (s UniqueSet) Add(value string) {
	s[value] = struct{}{}
}

// Contains checks if a string exists in the set
func (s UniqueSet) Contains(value string) bool {
	_, exists := s[value]
	return exists
}

func (s UniqueSet) Remove(value string) {
	delete(s, value)
}

func (s UniqueSet) Join(sep string) string {
	var result string
	for ele := range s {
		result += ele + sep
	}
	return result
}

/*
███████ ████████ ██████  ██ ███    ██  ██████  ███████
██         ██    ██   ██ ██ ████   ██ ██       ██
███████    ██    ██████  ██ ██ ██  ██ ██   ███ ███████
     ██    ██    ██   ██ ██ ██  ██ ██ ██    ██      ██
███████    ██    ██   ██ ██ ██   ████  ██████  ███████
*/

// remove empty strings, line breaks, and extra spaces
func CleanString(str string) string {
	// Regular expression to match multiple whitespace characters
	whitespace_regex := regexp.MustCompile(`\s+`)
	line_breakRegex := regexp.MustCompile(`(\r\n|\r|\n)`)

	// Replace all occurrences of the regex with a single space
	str = whitespace_regex.ReplaceAllString(str, " ")
	str = line_breakRegex.ReplaceAllString(str, "")
	return str

}

// split string at next given separator and return the two parts
func SplitAt(s, sep string) (string, string) {
	i := strings.Index(s, sep)
	if i == -1 {
		return s, ""
	}
	return s[:i], s[i+len(sep):]
}

func SplitAtFirst(s string, seps []string) (string, string) {
	for _, sep := range seps {
		i := strings.Index(s, sep)
		if i != -1 {
			return s[:i], s[i+len(sep):]
		}
	}
	return s, ""
}

func SplitAtRune(s string, r rune) (string, string) {
	i := strings.IndexRune(s, r)
	if i == -1 {
		return s, ""
	}
	return s[:i], s[i+1:]
}

func FindComponentTagEnd(s string, start int) int {
	inQuotes := false
	tagName := ""
	isClosingTag := false
	isSelfClosing := false

	// Find the tag name
	for i := start + 1; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '>' {
			tagName = s[start+1 : i]
			break
		}
	}

	for i := start; i < len(s); i++ {
		switch s[i] {
		case '"':
			inQuotes = !inQuotes
		case '/':
			if !inQuotes && i+1 < len(s) && s[i+1] == '>' {
				isSelfClosing = true
			}
		case '<':
			if !inQuotes && i+1 < len(s) && s[i+1] == '/' {
				isClosingTag = true
			}
		case '>':
			if !inQuotes {
				if isSelfClosing {
					return i
				}
				if isClosingTag {
					closingTagName := s[i-len(tagName)-2 : i-1]
					if closingTagName == tagName {
						return i
					}
				}
			}
		}
	}

	return -1 // Not found
}
