package tags

import (
	"fmt"
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/template/structure"
)

// DefineTag allows defining a named block of content that can be overridden by extending templates
var DefineTag = TemplateTag{
	Name: "define",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		// Extract the block name from the tag contents
		blockName := strings.TrimSpace(tag_contents)
		if blockName == "" {
			errs = append(errs, fmt.Errorf("define tag is missing the block name"))
			return pos, errs
		}

		// Find the closing tag
		closingTag := delimWrap(ctx, "enddefine")
		closingPos := strings.Index(raw[pos:], closingTag)
		if closingPos == -1 {
			errs = append(errs, fmt.Errorf("define tag missing closing enddefine"))
			return pos, errs
		}
		closingPos += pos

		// Get the content between the define tags
		content := raw[pos:closingPos]

		// Store the defined block in the context
		if ctx.Blocks == nil {
			ctx.Blocks = make(map[string]string)
		}
		ctx.Blocks[blockName] = content

		// Return position after the closing tag
		return closingPos + len(closingTag), errs
	},
}

// BlockTag represents a block that can be overridden by extending templates
var BlockTag = TemplateTag{
	Name: "block",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		// Extract the block name from the tag contents
		blockName := strings.TrimSpace(tag_contents)
		if blockName == "" {
			errs = append(errs, fmt.Errorf("block tag is missing the block name"))
			return pos, errs
		}

		// Find the closing tag
		closingTag := delimWrap(ctx, "endblock")
		closingPos := strings.Index(raw[pos:], closingTag)
		if closingPos == -1 {
			errs = append(errs, fmt.Errorf("block tag missing closing endblock"))
			return pos, errs
		}
		closingPos += pos

		// Get the default content between the block tags
		defaultContent := raw[pos:closingPos]

		// If this block is being extended, the extended content will be in ctx.Blocks
		if ctx.Blocks != nil {
			if extendedContent, exists := ctx.Blocks[blockName]; exists {
				// Render the extended content
				renderErrs := core.Render(ctx, sb, extendedContent)
				if len(renderErrs) > 0 {
					errs = append(errs, renderErrs...)
				}
				return closingPos + len(closingTag), errs
			}
		}

		// Otherwise render the default content
		renderErrs := core.Render(ctx, sb, defaultContent)
		if len(renderErrs) > 0 {
			errs = append(errs, renderErrs...)
		}

		return closingPos + len(closingTag), errs
	},
}
