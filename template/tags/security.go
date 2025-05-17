package tags

import (
	"fmt"
	"strings"

	"github.com/kato-studio/wispy/template/structure"
)

var RedirectTag = TemplateTag{
	Name: "redirect",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		// var errs []error

		// Parse tag options
		options := parseAssetTagOptions(tag_contents)
		redirect := strings.TrimSpace(options["redirect"])
		fmt.Print(redirect)

		return pos, nil
	},
}
