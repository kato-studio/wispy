package tags

import (
	"fmt"
	"maps"
	"reflect"
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/wispy_common/structure"
)

// Render
var EachTag = TemplateTag{
	Name: "each",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
		endTag := delimWrap(ctx, "end-each")
		openingOfTag := ctx.Engine.DelimStart + " each "
		endTagStart, endTagLength := core.SeekClosingHandleNested(raw, endTag, openingOfTag, pos)
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
		collection, varErr := core.ResolveVariable(ctx, collectionPath)
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
				newCtx := structure.RenderCtx{
					Engine:          ctx.Engine,
					Data:            newData, // Use cloned data
					Props:           maps.Clone(ctx.Props),
					ScopedDirectory: ctx.ScopedDirectory,
				}

				// Render block with new context
				var blockSB strings.Builder
				if renderErrs := core.Render(&newCtx, &blockSB, blockContent); len(renderErrs) > 0 {
					errs = append(errs, renderErrs...)
				}
				sb.WriteString(blockSB.String())
			}

		case reflect.Map:
			iter := collValue.MapRange()
			for iter.Next() {
				newData := maps.Clone(ctx.Data)
				newData[loopVar] = iter.Value().Interface() // Store in Data
				newCtx := structure.RenderCtx{
					Engine:          ctx.Engine,
					Data:            newData,
					Props:           maps.Clone(ctx.Props),
					ScopedDirectory: ctx.ScopedDirectory,
				}

				var blockSB strings.Builder
				if renderErrs := core.Render(&newCtx, &blockSB, blockContent); len(renderErrs) > 0 {
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
