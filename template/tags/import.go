package tags

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kato-studio/wispy/wispy_common/structure"
)

// relative import for adjacent files
var ImportTag = TemplateTag{
	Name: "import",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		var errs []error
		var options = parseAssetTagOptions(tag_contents)
		var path = options["path"]
		var external = options["external"] == "true"
		var isInline = options["inline"] == "true"
		var contentStr = ""

		if strings.HasPrefix(path, "https://") {
			external = true
		} else if options["inline"] != "true" {
			isInline = true
		}

		switch {
		case strings.HasPrefix(path, "~/"):
			external = false
			isInline = true
			// Relative to current template
			baseDir := filepath.Dir(ctx.CurrentTemplatePath)
			path = filepath.Join(baseDir, strings.TrimPrefix(path, "~/"))

			// Read file
			content, err := os.ReadFile(path)
			if err != nil {
				errs = append(errs, fmt.Errorf("import error: %v", err))
				return pos, errs
			}
			// Determine type and process
			contentStr = string(content)
		case external == false:
			baseDir := filepath.Dir("./")
			path = filepath.Join(baseDir, ctx.ScopedDirectory, path)

			// Read file
			content, err := os.ReadFile(path)
			if err != nil {
				errs = append(errs, fmt.Errorf("import error: %v", err))
				return pos, errs
			}
			// Determine type and process
			contentStr = string(content)
		}

		// get file type from extension
		ext := strings.ToLower(filepath.Ext(path))

		// Get priority if specified
		priority, _ := strconv.Atoi(options["priority"])

		switch ext {
		case ".css":
			asset := structure.Asset{
				Path:      path,
				Type:      structure.CSS,
				Content:   contentStr,
				IsInline:  isInline,
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
				Path:      path,
				Type:      structure.JS,
				Content:   contentStr,
				IsInline:  isInline,
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
			fmt.Println("could not determine asset type of, " + path)
			sb.WriteString(contentStr)
		}

		return pos, errs
	},
}
