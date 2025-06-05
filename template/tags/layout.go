package tags

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/wispy_common/structure"
)

// LayoutTag allows a template to specify a layout template that wraps its content
var LayoutTag = TemplateTag{
	Name: "layout",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		skipLayout := ctx.Request.Header.Get("HX-Skip-Layout")
		hxBoosted := ctx.Request.Header.Get("HX-Boosted")

		closingPos := len(raw)
		if strings.ToLower(skipLayout) == "true" && strings.ToLower(hxBoosted) != "true" {
			renderErrs := core.Render(ctx, sb, raw[pos:closingPos])
			if len(renderErrs) > 0 {
				errs = append(errs, renderErrs...)
			}

			return pos, errs
		}

		// Extract the layout template name from the tag contents
		layoutName := strings.Trim(tag_contents, " \"'")
		if layoutName == "" {
			errs = append(errs, fmt.Errorf("layout tag is missing the layout template name"))
			return pos, errs
		}

		// Get the content that will be wrapped by the layout
		content := raw[pos:closingPos]

		sitePartialsPath := filepath.Join(ctx.ScopedDirectory, "layouts")
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

		// Render contents pass to layout
		var tempBuilder strings.Builder
		renderErrs := core.Render(ctx, &tempBuilder, content)
		if len(renderErrs) > 0 {
			errs = append(errs, renderErrs...)
		}
		ctx.Passed = tempBuilder.String()

		// Update for use in asset imports
		// prevTemplatePath := ctx.CurrentTemplatePath
		ctx.CurrentTemplatePath = layoutFilePath

		// Render the layout template
		renderErrs = core.Render(ctx, sb, string(layoutContentAsBytes))
		if len(renderErrs) > 0 {
			errs = append(errs, renderErrs...)
		}

		// Return position after the closing tag
		return closingPos, errs
	},
}
