package tags

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/wispy_common/structure"
)

// PartialTag is a template tag that loads and renders a partial template.
// It expects the tag_contents to be the name of the partial file (without extension).
// The partial file should be located in the site's partials directory.
// The partial file is read, rendered, and the result is written to the string builder.
var PartialTag = TemplateTag{
	Name: "partial",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {

		// Extract the partial name from the tag contents
		partialName := strings.Trim(tag_contents, " \"'")
		if partialName == "" {
			errs = append(errs, fmt.Errorf("partial tag is missing the partial name"))
			return pos, errs
		}

		sitePartialsPath := filepath.Join(ctx.ScopedDirectory, "partials")
		partialFilePath := filepath.Join(sitePartialsPath, partialName+".hstm")

		partialContentAsBytes, err := os.ReadFile(partialFilePath)
		if err != nil {
			partialFilePath = filepath.Join(sitePartialsPath, partialName, "index.hstm")
			partialContentAsBytes, err = os.ReadFile(partialFilePath)
		}

		if err != nil {
			errs = append(errs, fmt.Errorf("failed to read layout template file %s: %v", partialFilePath, err))
			return pos, errs
		}

		var partialSB strings.Builder

		// Update for use in asset imports
		ctx.CurrentTemplatePath = partialFilePath

		renderErrs := core.Render(ctx, &partialSB, string(partialContentAsBytes))
		if len(renderErrs) > 0 {
			errs = append(errs, renderErrs...)
			return pos, errs
		}

		sb.WriteString(partialSB.String())

		// Return the new position (which is the same as the input pos since we didn't handle any content beyond the tag_contents)
		return pos, errs
	},
}
