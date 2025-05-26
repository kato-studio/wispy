package structure

import (
	"strings"
)

type HeadTag struct {
	TagName    string   // "meta", "link", etc.
	Attributes []string // "formatted tag attribute"
	Content    string   // For tags like title
}

type HeadTagRegistry struct {
	tags []*HeadTag
	seen map[string]struct{} // For deduplication
}

func (r *HeadTagRegistry) Add(tag *HeadTag) {
	// Create unique key for deduplication
	key := tag.TagName + strings.Join(tag.Attributes, ":")
	if _, exists := r.seen[key]; !exists {
		r.tags = append(r.tags, tag)
		r.seen[key] = struct{}{}
	}
}

func (r *HeadTagRegistry) Render() string {
	var sb strings.Builder
	for _, tag := range r.tags {
		switch tag.TagName {
		case "meta":
			sb.WriteString("<")
			sb.WriteString(tag.TagName)
			for _, attr := range tag.Attributes {
				sb.WriteString(" ")
				sb.WriteString(attr)
			}
			if tag.Content != "" {
				sb.WriteString(">")
				sb.WriteString(tag.Content)
				sb.WriteString("</")
				sb.WriteString(tag.TagName)
				sb.WriteString(">")
			} else {
				sb.WriteString(">")
			}
		default:
			sb.WriteString("<!-- UNKNOWN TAG: ")
			sb.WriteString(tag.TagName)
			sb.WriteString(" -->")
		}
	}

	return sb.String()
}
