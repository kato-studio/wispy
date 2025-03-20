package template

import (
	"fmt"
	"strings"
)

func ResolveFiltersIfAny(ctx *RenderCtx, sb *strings.Builder, tag_contents string) error {
	sb.WriteString("(" + tag_contents + ")")
	return nil
}

func FindDelim(ctx *RenderCtx, raw string, pos int) (int, int) {
	var ds = ctx.Engine.DelimStart
	var de = ctx.Engine.DelimEnd
	// Find the next occurrence of a variable or tag start delimiter.
	next := SafeIndex(raw, ds, pos)
	if next == -1 {
		return -1, -1
	}
	// find bracket closing delim
	endDelim := SafeIndex(raw, de, pos)
	return next, endDelim
}

func ResolveTag(ctx *RenderCtx, sb *strings.Builder, pos int, tag_contents, raw string) (new_pos int, errs []error) {
	tagName, contents, tagNameExists := strings.Cut(tag_contents, " ")
	if !tagNameExists {
		return pos, []error{fmt.Errorf("could not resolve tag name in \"" + tag_contents + "\"")}
	}

	// check if tag has been registered to the template engine.
	templateTag, tagExists := ctx.Engine.TagMap[tagName]
	if !tagExists {
		return pos, []error{fmt.Errorf("\"" + tagName + "\" not found in Engine.TagMap \"" + tag_contents + "\"")}
	}

	return templateTag.Render(ctx, sb, contents, raw, pos)
}

// resolveVariable resolves a variable reference from the RenderCtx's Props or Data maps.
func ResolveVariable(ctx *RenderCtx, path string) (any, error) {
	parts := strings.Split(path, ".")
	var current any
	// try to resolve the variable from Props.
	current = ctx.Props
	for _, part := range parts {
		if part == "" {
			continue // Skip empty parts (e.g., leading dot).
		}
		// Check if the current value is a map.
		if m, ok := current.(map[string]any); ok {
			if val, exists := m[part]; exists {
				current = val
			} else {
				// Variable not found in Props, fall back to Data.
				current = nil
				break
			}
		} else {
			// Invalid path in Props, fall back to Data.
			current = nil
			break
		}
	}
	// If the variable was not found in Props, try to resolve it from Data.
	if current == nil {
		current = ctx.Data
		for _, part := range parts {
			if part == "" {
				continue // Skip empty parts (e.g., leading dot).
			}
			// Check if the current value is a map.
			if m, ok := current.(map[string]any); ok {
				if val, exists := m[part]; exists {
					current = val
				} else {
					// Variable not found in Data.
					fmt.Printf("variable not found: %s\n", part)
					return "", nil
				}
			} else {
				// Invalid path in Data.
				return "", fmt.Errorf("invalid variable path: %s", path)
			}
		}
	}
	return current, nil
}

// evaluateCondition evaluates a condition string (e.g., "x > 5").
func ResolveCondition(ctx *RenderCtx, condition string) (val bool, errs []error) {
	// Simple (for demonstration purposes)
	// TODO:  condition evaluation I.E. {{ if eq 1 1 }}
	return
}
