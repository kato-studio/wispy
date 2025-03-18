package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Universal template tag function struct
type TemplateTag struct {
	Name string
	// render tag with given context and args, and children if a RequiresClosingTag is set
	Render func(
		// The context of this render execution including...
		// - Reference to the template engine struct
		// - Partials map,
		// - Data map fetched via eng.GetFunc(ctx *RenderCtx, key string)
		ctx *RenderCtx,
		// Finalized output is written to the string building
		sb *strings.Builder,
		// The inner contexts of the tag being parsed
		// Example: "{% exampleTag ... ... ... %}"
		// (tags using closing tag for inner content are expected to resolve the closing tag and content then return update POS int)
		tag_contents,
		// entire input string mainly used if tag handles closing child_content + closing_tag
		raw string,
		// the current position within the raw input string
		pos int,
	) (new_pos int, errs []error)
}

// PartialTag is a template tag that loads and renders a partial template.
// It expects the tag_contents to be the name of the partial file (without extension).
// The partial file should be located in the site's partials directory.
// The partial file is read, rendered, and the result is written to the string builder.
var PartialTag = TemplateTag{
	Name: "partial",
	Render: func(ctx *RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		// Extract the partial name from the tag contents
		partialName := strings.Trim(tag_contents, " \"'")
		if partialName == "" {
			errs = append(errs, fmt.Errorf("partial tag is missing the partial name"))
			return pos, errs
		}

		// Construct the path to the partial file
		sitePartialsPath := filepath.Join(ctx.ScopedDirectory, "partials")
		partialFilePath := filepath.Join(sitePartialsPath, partialName+".hstm")

		// Read the partial file
		partialContentAsBytes, err := os.ReadFile(partialFilePath)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to read partial file %s: %v", partialFilePath, err))
			return pos, errs
		}

		// Render the partial content
		// Note: This assumes that the partial content is a template that needs to be rendered.
		// You might need to recursively render the partial content if it contains other tags.
		var partialSB strings.Builder
		renderErrs := Render(ctx, &partialSB, string(partialContentAsBytes))
		if len(renderErrs) > 0 {
			errs = append(errs, renderErrs...)
			return pos, errs
		}

		// Write the rendered partial to the string builder
		sb.WriteString(partialSB.String())

		// Return the new position (which is the same as the input pos since we didn't move in the raw string)
		return pos, errs
	},
}

var IfTag = TemplateTag{
	Name: "if",
	Render: func(ctx *RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		endTag := IndexAt(raw, "{% endif %}", pos)
		if endTag == -1 {
			errs = append(errs, fmt.Errorf("could not find end tag for %s", "{% endtag %}"))
			return pos, errs
		}

		sb.WriteString("[[" + tag_contents + "]]")
		return pos, errs
	},
}

var ForTag = TemplateTag{
	Name: "for",
	Render: func(ctx *RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		sb.WriteString("[[" + tag_contents + "]]")
		return pos, errs
	},
}

var RootTag = TemplateTag{
	Name: "root",
	Render: func(ctx *RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		return pos, errs
	},
}

var DefaultTemplateTags = []TemplateTag{
	IfTag,
	ForTag,
	PartialTag,
}

var DefaultEngineTags = []TemplateTag{
	IfTag,
	ForTag,
	PartialTag,
	// PageTag,
	RootTag,
}
