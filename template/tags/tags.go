package tags

import (
	"fmt"
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/wispy_common"
	"github.com/kato-studio/wispy/wispy_common/structure"
)

type TemplateTag = structure.TemplateTag

func delimWrap(ctx *structure.RenderCtx, value string) string {
	return strings.Join([]string{ctx.Engine.DelimStart, value, ctx.Engine.DelimEnd}, " ")
}

// parseAssetTagOptions handles the full parsing including the path edge case
// and guarantees path variable is return
func parseAssetTagOptions(input string) map[string]string {
	pairs := core.SplitRespectQuotes(input)
	options := wispy_common.ParseKeyValuePairs(pairs)

	// Handle the case for standalone path
	if _, exists := options["path"]; !exists {
		if len(pairs) > 0 {
			firstPair := strings.Trim(pairs[0], `"'`)
			options["path"] = firstPair
		} else {
			options["path"] = "no-path-supplied"
		}

	}

	return options
}

// insertNestedValue uses maps.Insert to set values in nested structures
func insertNestedValue(root map[string]interface{}, path []string, value interface{}) error {
	if len(path) == 0 {
		return fmt.Errorf("empty path")
	}

	// For single segment paths, just do direct assignment
	if len(path) == 1 {
		root[path[0]] = value
		return nil
	}

	// Create the nested structure if needed
	current := root
	for i := 0; i < len(path)-1; i++ {
		part := path[i]
		if next, exists := current[part]; exists {
			if nextMap, ok := next.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return fmt.Errorf("path segment '%s' is not a map", part)
			}
		} else {
			newMap := make(map[string]interface{})
			current[part] = newMap
			current = newMap
		}
	}

	// Set the final value
	current[path[len(path)-1]] = value
	return nil
}

// used to skip content that should not be parsed
var CommentTag = TemplateTag{
	Name: "comment",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		endTag := delimWrap(ctx, "end-comment")
		endTagStart, endTagLength := core.SeekIndexAndLength(raw, endTag, pos)
		if endTagStart == -1 {
			errs = append(errs, fmt.Errorf("could not find end tag for %s", endTag))
			return pos, errs
		}
		newEndPos := endTagStart + endTagLength

		return newEndPos, errs
	},
}
