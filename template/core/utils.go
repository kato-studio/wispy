package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kato-studio/wispy/template/structure"
	"github.com/kato-studio/wispy/utilities"
)

// returns the index of the from current pos, if non-found the pos will be the same upon return
func SeekIndex(raw, sep string, pos int) (new_pos int) {
	new_pos = strings.Index(raw[pos:], sep)
	if new_pos > -1 {
		new_pos += pos
	}
	return new_pos
}

func SeekIndexAndLength(raw, sep string, pos int) (new_pos int, seperator_lenth int) {
	new_pos = strings.Index(raw[pos:], sep)
	if new_pos > -1 {
		new_pos += pos
	}
	return new_pos, len(sep)
}

// SafeIndexAndLength finds the index of a closing tag while ensuring it corresponds to an opening tag
func SeekClosingHandleNested(raw, closingTag, openingTag string, pos int) (newPos int, separatorLength int) {
	openCount := 0
	newPos = pos
	separatorLength = len(closingTag)

	for {
		closeIndex := strings.Index(raw[newPos:], closingTag)
		if closeIndex == -1 {
			// No more closing tags found, return -1
			return -1, separatorLength
		}
		closeIndex += newPos

		// Only search for an opening tag within the range before this closing tag
		openIndex := strings.Index(raw[newPos:closeIndex], openingTag)
		if openIndex != -1 {
			openCount++
			newPos += openIndex + len(openingTag)
			continue
		}

		// If no unmatched opening tags remain, return this closing tag index
		if openCount == 0 {
			return closeIndex, separatorLength
		}

		// Otherwise, decrement open count and continue searching
		openCount--
		newPos = closeIndex + separatorLength
	}
}

func ResolveValue(ctx *structure.RenderCtx, expr string) (any, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, fmt.Errorf("empty expression")
	}

	// Check if it's a variable (starts with .)
	if strings.HasPrefix(expr, ".") {
		return ResolveVariable(ctx, expr)
	}

	// Check if it's a string literal
	if (strings.HasPrefix(expr, `"`) && strings.HasSuffix(expr, `"`)) ||
		(strings.HasPrefix(expr, `'`) && strings.HasSuffix(expr, `'`)) {
		return expr[1 : len(expr)-1], nil
	}

	// Check if it's a boolean
	if expr == "true" {
		return true, nil
	}
	if expr == "false" {
		return false, nil
	}

	// Check if it's a numeric value
	if num, err := strconv.Atoi(expr); err == nil {
		return num, nil
	}
	if num, err := strconv.ParseFloat(expr, 64); err == nil {
		return num, nil
	}

	// Check if it's a variable without . prefix (assuming it's a global or other scope)
	return ResolveVariable(ctx, expr)
}

// ---
// The Following functions have been moved to a general utilities package
// ---

// SplitRespectQuotes splits a string by spaces while respecting quoted substrings and removing empty values.
var SplitRespectQuotes = utilities.SplitRespectQuotes

var ParseDataPath = utilities.ParseDataPath

// this stringify function currently handles non-string values when resolving variables from ctx.Data & ctx.Props
var stringify = utilities.stringify
