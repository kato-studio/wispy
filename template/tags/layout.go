package tags

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/template/structure"
)

// LayoutTag allows a template to specify a layout template that wraps its content
var LayoutTag = TemplateTag{
	Name: "layout",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		// Extract the layout template name from the tag contents
		layoutName := strings.Trim(tag_contents, " \"'")
		if layoutName == "" {
			errs = append(errs, fmt.Errorf("layout tag is missing the layout template name"))
			return pos, errs
		}

		closingPos := len(raw)

		// Get the content that will be wrapped by the layout
		content := raw[pos:closingPos]

		// Read the layout template file
		sitePartialsPath := filepath.Join(ctx.ScopedDirectory, "partials")
		layoutFilePath := filepath.Join(sitePartialsPath, layoutName+".hstm")

		layoutContentAsBytes, err := os.ReadFile(layoutFilePath)
		if err != nil {
			layoutFilePath = filepath.Join(sitePartialsPath, layoutName, "index.hstm")
			layoutContentAsBytes, err = os.ReadFile(layoutFilePath)
		}

		if err != nil {
			errs = append(errs, fmt.Errorf("failed to read layout template file %s: %v", layoutFilePath, err))
			return pos, errs
		}

		// Update for use in asset imports
		ctx.CurrentTemplatePath = layoutFilePath

		// Render contents passed to layout
		var tempBuilder strings.Builder
		renderErrs := core.Render(ctx, &tempBuilder, content)
		if len(renderErrs) > 0 {
			errs = append(errs, renderErrs...)
		}
		ctx.Passed = tempBuilder.String()

		// Render the layout template
		renderErrs = core.Render(ctx, sb, string(layoutContentAsBytes))
		if len(renderErrs) > 0 {
			errs = append(errs, renderErrs...)
		}

		// Return position after the closing tag
		return closingPos, errs
	},
}
