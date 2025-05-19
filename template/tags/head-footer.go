package tags

import (
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/wispy_common"
	"github.com/kato-studio/wispy/wispy_common/structure"
)

// parseAssetTagOptions handles the full parsing including the path edge case
func parseAssetTagOptions(input string) map[string]string {
	pairs := core.SplitRespectQuotes(input)
	options := wispy_common.ParseKeyValuePairs(pairs)

	// Handle the special case for standalone path
	if _, exists := options["path"]; !exists && len(pairs) > 0 {
		firstPair := strings.Trim(pairs[0], `"'`)
		if !strings.Contains(firstPair, "=") {
			options["path"] = firstPair
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

		tag := structure.HeadTag{
			TagName:    "link",
			Attributes: options,
		}

		ctx.HeadTags.Add(&tag)
		return pos, nil
	},
}

var MetaTag = TemplateTag{
	Name: "meta",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		options := core.SplitRespectQuotes(tag_contents)

		tag := structure.HeadTag{
			TagName:    "meta",
			Attributes: options,
		}

		ctx.HeadTags.Add(&tag)
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
