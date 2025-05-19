package tags

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kato-studio/wispy/wispy_common/structure"
)

var ImportCSSTag = TemplateTag{
	Name: "import-css",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		var errs []error
		options := parseAssetTagOptions(tag_contents)

		// Handle both inline and external CSS
		if options["inline"] != "" && options["path"] != "" {
			errs = append(errs, fmt.Errorf("cannot specify both path and inline content"))
			return pos, errs
		}

		priority, err := strconv.Atoi(options["priority"])
		if err != nil {
			priority = 100 // Default CSS priority
		}

		asset := structure.Asset{
			Path:      options["path"],
			Type:      structure.CSS,
			Priority:  priority,
			Content:   options["inline"],
			IsInline:  options["inline"] != "",
			Media:     options["media"],             // Support media queries
			Preload:   options["preload"] == "true", // CSS preloading
			Condition: options["if"],                // Conditional loading
		}

		// Validate before adding
		if !asset.IsInline && asset.Path == "" {
			errs = append(errs, fmt.Errorf("external CSS requires a path"))
			return pos, errs
		}

		ctx.AssetRegistry.Add(&asset)
		return pos, errs
	},
}

var ImportJSTag = TemplateTag{
	Name: "import-js",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		var errs []error
		options := parseAssetTagOptions(tag_contents)

		// Handle both inline and external JS
		if options["inline"] != "" && options["path"] != "" {
			errs = append(errs, fmt.Errorf("cannot specify both path and inline content"))
			return pos, errs
		}

		priority, err := strconv.Atoi(options["priority"])
		if err != nil {
			priority = 200 // Default JS priority
		}

		asset := structure.Asset{
			Path:        options["path"],
			Type:        structure.JS,
			Priority:    priority,
			Async:       options["async"] == "true",
			Defer:       options["defer"] == "true",
			Module:      options["module"] == "true",
			Nomodule:    options["nomodule"] == "true", // Legacy browser support
			Content:     options["inline"],
			IsInline:    options["inline"] != "",
			Integrity:   options["integrity"],   // SRI hash
			Crossorigin: options["crossorigin"], // CORS setting
			Condition:   options["if"],          // Conditional loading
		}

		// Validate before adding
		if !asset.IsInline && asset.Path == "" {
			errs = append(errs, fmt.Errorf("external JS requires a path"))
			return pos, errs
		}

		ctx.AssetRegistry.Add(&asset)
		return pos, errs
	},
}

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

// relative import for adjacent files
var ImportTag = TemplateTag{
	Name: "import",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		var errs []error
		options := parseAssetTagOptions(tag_contents)

		path := options["path"]
		if path == "" {
			// Support both {% import "path" %} and {% import path="path" %}
			path = strings.TrimSpace(tag_contents)
		}

		// Path resolution
		switch {
		case strings.HasPrefix(path, "~/"):
			// Relative to current template
			baseDir := filepath.Dir(ctx.CurrentTemplatePath)
			path = filepath.Join(baseDir, strings.TrimPrefix(path, "~/"))
			// case strings.HasPrefix(path, "@"):
			// 	// Relative to project root or perhaps aliases?
			// 	path = filepath.Join(ctx.ProjectRoot, strings.TrimPrefix(path, "@/"))
		}

		// Read file
		content, err := os.ReadFile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("import error: %v", err))
			return pos, errs
		}

		// Determine type and process
		ext := strings.ToLower(filepath.Ext(path))
		contentStr := string(content)

		// Get priority if specified
		priority, _ := strconv.Atoi(options["priority"])

		switch ext {
		case ".css":
			asset := structure.Asset{
				Type:      structure.CSS,
				Content:   contentStr,
				IsInline:  true,
				Priority:  priority,
				Media:     options["media"],
				Condition: options["if"],
			}
			if priority == 0 {
				asset.Priority = 100 // Default CSS
			}
			ctx.AssetRegistry.Add(&asset)

		case ".js":
			asset := structure.Asset{
				Type:      structure.JS,
				Content:   contentStr,
				IsInline:  true,
				Priority:  priority,
				Async:     options["async"] == "true",
				Defer:     options["defer"] == "true",
				Module:    options["module"] == "true",
				Condition: options["if"],
			}
			if priority == 0 {
				asset.Priority = 200 // Default JS
			}
			ctx.AssetRegistry.Add(&asset)

		default:
			// Directly include other file types
			sb.WriteString(contentStr)
		}

		return pos, errs
	},
}
