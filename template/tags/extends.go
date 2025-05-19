package tags

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/wispy_common/structure"
)

// ExtendsTag allows a template to extend another template and override its blocks
var ExtendsTag = TemplateTag{
	Name: "extends",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		// Extract the parent template name from the tag contents
		parentName := strings.Trim(tag_contents, " \"'")
		if parentName == "" {
			errs = append(errs, fmt.Errorf("extends tag is missing the parent template name"))
			return pos, errs
		}

		// Find the closing tag for the entire content (if any)
		// This is optional as extends might be at the top with content following
		closingTag := delimWrap(ctx, "end-extends")
		closingPos := strings.Index(raw[pos:], closingTag)
		if closingPos != -1 {
			closingPos += pos
		} else {
			closingPos = len(raw)
		}

		// Get the child content (blocks to override)
		childContent := raw[pos:closingPos]

		// Read the parent template
		sitePartialsPath := filepath.Join(ctx.ScopedDirectory, "partials")
		parentFilePath := filepath.Join(sitePartialsPath, parentName+".hstm")

		parentContentAsBytes, err := os.ReadFile(parentFilePath)
		if err != nil {
			parentFilePath = filepath.Join(sitePartialsPath, parentName, "index.hstm")
			parentContentAsBytes, err = os.ReadFile(parentFilePath)
		}

		if err != nil {
			errs = append(errs, fmt.Errorf("failed to read layout template file %s: %v", parentFilePath, err))
			return pos, errs
		}

		// Update for use in asset imports
		ctx.CurrentTemplatePath = parentFilePath

		// Parse the child content to find block definitions
		slotStartTag := delimWrap(ctx, "slot ")
		slotEndTag := delimWrap(ctx, "end-slot")

		var rawContent strings.Builder
		childPos := 0
		for childPos < len(childContent) {
			startPos := strings.Index(childContent[childPos:], slotStartTag)
			if startPos == -1 {
				// catch content if no slot tags or trailing content after all slots have been captured
				rawContent.WriteString(childContent[childPos:])
				break
			}
			// handle any content between slot tags
			rawContent.WriteString(childContent[childPos:startPos])

			// update active pos
			startPos += childPos

			endPos := strings.Index(childContent[startPos:], slotEndTag)
			if endPos == -1 {
				errs = append(errs, fmt.Errorf("unclosed slot block in child template"))
				break
			}
			endPos += startPos

			// Extract slot name (after "define ")
			slotNameStart := startPos + len(slotStartTag)
			slotNameEnd := strings.IndexAny(childContent[slotNameStart:], " \t\n\r")
			if slotNameEnd == -1 {
				slotNameEnd = endPos
			} else {
				slotNameEnd += slotNameStart
			}

			slotName := strings.TrimSpace(childContent[slotNameStart:slotNameEnd])
			slotContent := childContent[slotNameEnd:endPos]

			// Store the block content
			ctx.Slots[slotName] = slotContent

			childPos = endPos + len(slotEndTag)
		}

		var renderErrs []error
		// Render contents passed to layout
		var tempBuilder strings.Builder
		pased := strings.TrimSpace(rawContent.String())
		if len(pased) > 0 {
			renderErrs = core.Render(ctx, &tempBuilder, pased)
			if len(renderErrs) > 0 {
				errs = append(errs, renderErrs...)
			}
			ctx.Passed = tempBuilder.String()
		}

		// Render the parent template with the child blocks
		renderErrs = core.Render(ctx, sb, string(parentContentAsBytes))
		if len(renderErrs) > 0 {
			errs = append(errs, renderErrs...)
		}

		// Return position after the closing tag (or end of input)
		if closingPos == len(raw) {
			return closingPos, errs
		}
		return closingPos + len(closingTag), errs
	},
}
