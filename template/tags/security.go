package tags

var RedirectTag = TemplateTag{
	Name: "redirect",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
			var errs []error
			
			// Parse tag options
			options := parseAssetTagOptions(tag_contents)
			redirect := strings.TrimSpace(options["redirect"])
			
			// Find the end tag
			endTag := delimWrap(ctx, "endprotected")
			endTagStart, endTagLength := core.SeekIndexAndLength(raw, endTag, pos)
			if endTagStart == -1 {
					errs = append(errs, fmt.Errorf("could not find end tag for %s", endTag))
					return pos, errs
			}
		}
	}
	