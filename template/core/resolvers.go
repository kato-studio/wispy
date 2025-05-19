package core

import (
	"fmt"
	"maps"
	"strings"

	"github.com/kato-studio/wispy/wispy_common/structure"
)

func ResolveFiltersIfAny(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents string) error {
	// sb.WriteString("(" + tag_contents + ")")
	parts := strings.Split(tag_contents, " | ")
	lenParts := len(parts)
	//
	val, valErr := ResolveVariable(ctx, parts[0])
	if valErr != nil {
		return valErr
	}
	if lenParts == 1 {
		if str, ok := val.(string); ok {
			sb.WriteString(str)
		} else {
			sb.WriteString(Stringify(val))
		}
	} else if lenParts > 1 {
		var pipeValue any = val
		var err error
		for i := 1; i < lenParts; i++ {
			filterSec := parts[i]
			filterParts := SplitRespectQuotes(filterSec)
			filterName := filterParts[1]
			//
			if filter, ok := ctx.Engine.FilterMap[filterName]; ok {
				pipeValue, err = filter.Handler(pipeValue, filterParts)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("no filter found %s", filterName)
			}
		}
	}
	return nil
}

func FindDelim(ctx *structure.RenderCtx, raw string, pos int) (int, int) {
	var ds = ctx.Engine.DelimStart
	var de = ctx.Engine.DelimEnd
	// Find the next occurrence of a variable or tag start delimiter.
	next := SeekIndex(raw, ds, pos)
	if next == -1 {
		return -1, -1
	}
	// find bracket closing delim
	endDelim := SeekIndex(raw, de, pos)
	return next, endDelim
}

func ResolveTag(ctx *structure.RenderCtx, sb *strings.Builder, pos int, tag_contents, raw string) (new_pos int, errs []error) {
	tagName, contents, tagNameExists := strings.Cut(tag_contents, " ")
	if !tagNameExists && len(tagName) < 3 {
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
func ResolveVariable(ctx *structure.RenderCtx, path string) (any, error) {
	parts := strings.Split(path, ".")
	var current any = maps.Clone(ctx.Props)
	// try to resolve the variable from Props.
	current = ParseDataPath(parts, current)
	// If the variable was not found in Props, try to resolve it from Data.
	if current == nil {
		current = ctx.Data
		current = ParseDataPath(parts, current)
	}
	// If it's still we know there was no data either so we return an error
	if current == nil {
		return "", fmt.Errorf("invalid variable path: %s", path)
	}
	return current, nil
}
