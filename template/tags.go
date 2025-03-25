package template

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"reflect"
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

		sitePartialsPath := filepath.Join(ctx.ScopedDirectory, "partials")
		partialFilePath := filepath.Join(sitePartialsPath, partialName+".hstm")

		partialContentAsBytes, err := os.ReadFile(partialFilePath)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to read partial file %s: %v", partialFilePath, err))
			return pos, errs
		}

		var partialSB strings.Builder
		renderErrs := Render(ctx, &partialSB, string(partialContentAsBytes))
		if len(renderErrs) > 0 {
			errs = append(errs, renderErrs...)
			return pos, errs
		}

		sb.WriteString(partialSB.String())

		// Return the new position (which is the same as the input pos since we didn't handle any content beyond the tag_contents)
		return pos, errs
	},
}

// If take only show content if a value is true
var IfTag = TemplateTag{
	Name: "if",
	Render: func(ctx *RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		endTagStart, endTagLength := SeekIndexAndLenth(raw, "{% endif %}", pos)
		if endTagStart == -1 {
			errs = append(errs, fmt.Errorf("could not find end tag for %s", "{% endif %}"))
			return pos, errs
		}
		content := raw[pos:endTagStart]
		newEndPos := endTagStart + endTagLength

		value, condition_errors := ResolveCondition(ctx, tag_contents)
		if len(condition_errors) > 0 {
			errs = append(errs, condition_errors...)
		}
		if value {
			sb.WriteString(content)
		}
		return newEndPos, errs
	},
}

// Render
var EachTag = TemplateTag{
	Name: "each",
	Render: func(ctx *RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		endTag := ctx.Engine.DelimStart + " endeach " + ctx.Engine.DelimEnd
		openingOfTag := ctx.Engine.DelimStart + " each "
		endTagStart, endTagLength := SeedClosingHandleNested(raw, endTag, openingOfTag, pos)
		if endTagStart == -1 {
			errs = append(errs, fmt.Errorf("could not find end tag %q", endTag))
			return pos, errs
		}
		newEndPos := endTagStart + endTagLength
		//
		// Parse loop variables and collection path
		parts := strings.SplitN(tag_contents, " in ", 2)
		if len(parts) != 2 {
			errs = append(errs, fmt.Errorf("invalid each syntax: expected '{VAR} in {ARRAY}', got %q", tag_contents))
			return newEndPos, errs
		}
		//
		loopVar := strings.TrimSpace(parts[0])
		//
		collectionPath := strings.TrimSpace(parts[1])
		//
		collection, varErr := ResolveVariable(ctx, collectionPath)
		if varErr != nil {
			errs = append(errs, varErr)
		}

		collValue := reflect.ValueOf(collection)
		blockContent := raw[pos:endTagStart]

		switch collValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < collValue.Len(); i++ {
				// Clone parent data and add loop variable
				newData := maps.Clone(ctx.Data)
				newData[loopVar] = collValue.Index(i).Interface() // Store in Data

				newCtx := &RenderCtx{
					Engine:          ctx.Engine,
					Data:            newData, // Use cloned data
					Props:           maps.Clone(ctx.Props),
					ScopedDirectory: ctx.ScopedDirectory,
				}

				// Render block with new context
				var blockSB strings.Builder
				if renderErrs := Render(newCtx, &blockSB, blockContent); len(renderErrs) > 0 {
					errs = append(errs, renderErrs...)
				}
				sb.WriteString(blockSB.String())
			}

		case reflect.Map:
			iter := collValue.MapRange()
			for iter.Next() {
				newData := maps.Clone(ctx.Data)
				newData[loopVar] = iter.Value().Interface() // Store in Data

				newCtx := &RenderCtx{
					Engine:          ctx.Engine,
					Data:            newData,
					Props:           maps.Clone(ctx.Props),
					ScopedDirectory: ctx.ScopedDirectory,
				}

				var blockSB strings.Builder
				if renderErrs := Render(newCtx, &blockSB, blockContent); len(renderErrs) > 0 {
					errs = append(errs, renderErrs...)
				}
				sb.WriteString(blockSB.String())
			}

		default:
			errs = append(errs, fmt.Errorf("cannot iterate over %T", collection))
		}

		return endTagStart + endTagLength, errs
	},
}

var DefaultTemplateTags = []TemplateTag{
	IfTag,
	EachTag,
	PartialTag,
}

var DefaultEngineTags = []TemplateTag{
	IfTag,
	EachTag,
	PartialTag,
}
