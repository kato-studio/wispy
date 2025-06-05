package tags

import (
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/wispy_common/structure"
)

var HeadTag = TemplateTag{
	Name: "root-head",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, _, _ string, pos int) (int, []error) {
		sb.WriteString(ctx.HeadTags.Render())
		return pos, nil
	},
}

var CssAssetsTag = TemplateTag{
	Name: "root-css",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, _, _ string, pos int) (int, []error) {
		sb.WriteString(ctx.AssetRegistry.Render(structure.CSS))
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

var JsAssetsTag = TemplateTag{
	Name: "root-js",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, _, _ string, pos int) (int, []error) {
		sb.WriteString(ctx.AssetRegistry.Render(structure.JS))
		return pos, nil
	},
}
