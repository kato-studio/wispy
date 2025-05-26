package tags

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kato-studio/wispy/wispy_common/structure"
)

// Block-style CSS tag
var CSSTag = TemplateTag{
	Name: "css",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		options := parseAssetTagOptions(tag_contents)

		// Find closing tag
		closingTag := delimWrap(ctx, "end-css")
		closingPos := strings.Index(raw[pos:], closingTag)
		if closingPos == -1 {
			return pos, []error{fmt.Errorf("missing closing endcss tag")}
		}
		closingPos += pos

		// Get content between tags
		content := strings.TrimSpace(raw[pos:closingPos])
		content = strings.TrimPrefix(content, "<style>")
		content = strings.TrimSuffix(content, "</style>")

		priority, _ := strconv.Atoi(options["priority"])

		ctx.AssetRegistry.Add(&structure.Asset{
			Type:     structure.CSS,
			Content:  content,
			IsInline: true,
			Priority: priority,
			Media:    options["media"],
		})

		return closingPos + len(closingTag), nil
	},
}

// Block-style JS tag
var JSTag = TemplateTag{
	Name: "js",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		options := parseAssetTagOptions(tag_contents)

		closingTag := delimWrap(ctx, "end-js")
		closingPos := strings.Index(raw[pos:], closingTag)
		if closingPos == -1 {
			return pos, []error{fmt.Errorf("missing closing endjs tag")}
		}
		closingPos += pos

		content := strings.TrimSpace(raw[pos:closingPos])
		content = strings.TrimPrefix(content, "<script>")
		content = strings.TrimSuffix(content, "</script>")

		ctx.AssetRegistry.Add(&structure.Asset{
			Type:     structure.JS,
			Content:  content,
			IsInline: true,
			Async:    options["async"] == "true",
			Defer:    options["defer"] == "true",
			Module:   options["module"] == "true",
		})

		return closingPos + len(closingTag), nil
	},
}
