package tags

import (
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/template/structure"
)

// this function is also used by /tags/css-js.go
func parseAssetTagOptions(input string) map[string]string {
	options := make(map[string]string)

	// Split into key=value pairs
	pairs := core.SplitRespectQuotes(input)
	for _, pair := range pairs {
		// Handle inline content (wrapped in quotes)
		if strings.Contains(pair, "=") {
			parts := strings.SplitN(pair, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			options[key] = value
		} else {
			// Handle standalone flags or default path
			if _, exists := options["path"]; !exists {
				options["path"] = strings.Trim(pair, `"'`)
			} else {
				options[pair] = "true"
			}
		}
	}

	return options
}

var HeadTag = TemplateTag{
	Name: "head",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, _, _ string, pos int) (int, []error) {
		sb.WriteString("<head>")
		sb.WriteString(ctx.HeadTags.Render())
		sb.WriteString(ctx.AssetRegistry.Render(structure.CSS))
		// sb.WriteString(ctx.AssetRegistry.Render(structure.JS))
		sb.WriteString("</head>")
		return pos, nil
	},
}

var TitleTag = TemplateTag{
	Name: "title",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		ctx.HeadTags.Add(&structure.HeadTag{
			TagName: "title",
			Content: strings.Trim(tag_contents, "\""),
		})
		return pos, nil
	},
}

var LinkTag = TemplateTag{
	Name: "link",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		options := core.SplitRespectQuotes(tag_contents)

		tag := &structure.HeadTag{
			TagName:    "link",
			Attributes: options,
		}

		ctx.HeadTags.Add(tag)
		return pos, nil
	},
}

var MetaTag = TemplateTag{
	Name: "meta",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		options := core.SplitRespectQuotes(tag_contents)

		tag := &structure.HeadTag{
			TagName:    "meta",
			Attributes: options,
		}

		ctx.HeadTags.Add(tag)
		return pos, nil
	},
}

var FooterAssetsTag = TemplateTag{
	Name: "footer",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, _, _ string, pos int) (int, []error) {
		// sb.WriteString("<foter>")
		sb.WriteString(ctx.AssetRegistry.Render(structure.JS))
		// sb.WriteString("</foter>")
		return pos, nil
	},
}
