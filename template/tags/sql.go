package tags

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/template/structure"

	_ "github.com/tursodatabase/go-libsql"
)

var SQLiteTag = TemplateTag{
	Name: "sqlite",
	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
		var errs []error

		// Parse tag options
		options := parseAssetTagOptions(tag_contents)
		query := strings.TrimSpace(options["query"])
		dbPath := strings.TrimSpace(options["path"])

		// Handle path resolution (similar to import tag)
		if strings.HasPrefix(dbPath, "~/") {
			baseDir := filepath.Dir(ctx.CurrentTemplatePath)
			dbPath = filepath.Join(baseDir, strings.TrimPrefix(dbPath, "~/"))
		} else if strings.HasPrefix(dbPath, "@root/") {
			// Assuming @root is a predefined root directory
			dbPath = filepath.Join(ctx.ScopedDirectory, strings.TrimPrefix(dbPath, "@root/"))
		}

		// Validate required parameters
		if query == "" {
			errs = append(errs, fmt.Errorf("sqlite tag requires a query parameter"))
			return pos, errs
		}
		if dbPath == "" {
			errs = append(errs, fmt.Errorf("sqlite tag requires a path parameter"))
			return pos, errs
		}

		// Find the end tag
		endTag := delimWrap(ctx, "endsqlite")
		endTagStart, endTagLength := core.SeekIndexAndLength(raw, endTag, pos)
		if endTagStart == -1 {
			errs = append(errs, fmt.Errorf("could not find end tag for %s", endTag))
			return pos, errs
		}

		// Extract content between tags
		content := raw[pos:endTagStart]
		newEndPos := endTagStart + endTagLength

		// Open database connection
		db, err := sql.Open("libsql", "file:"+dbPath)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to open database at %s: %v", dbPath, err))
			return newEndPos, errs
		}
		defer db.Close()

		// Prepare the query
		stmt, err := db.Prepare(query)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to prepare query: %v", err))
			return newEndPos, errs
		}
		defer stmt.Close()

		// Extract query parameters from template variables
		var params []interface{}
		// paramRegex := regexp.MustCompile(`\{\{\s*(.+?)\s*\}\}`)
		// matches := paramRegex.FindAllStringSubmatch(query, -1)

		// for _, match := range matches {
		// 	if len(match) > 1 {
		// 		val, err := core.ResolveValue(ctx, match[1])
		// 		if err != nil {
		// 			errs = append(errs, fmt.Errorf("failed to resolve parameter %s: %v", match[1], err))
		// 		} else {
		// 			params = append(params, val)
		// 		}
		// 	}
		// }

		// Clean the query by removing template variables
		// cleanQuery := paramRegex.ReplaceAllString(query, "?")

		// Execute the query
		rows, err := stmt.Query(params...)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to execute query: %v", err))
			return newEndPos, errs
		}
		defer rows.Close()

		// Get column names
		columns, err := rows.Columns()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get column names: %v", err))
			return newEndPos, errs
		}

		// Process query results
		var results []map[string]interface{}
		for rows.Next() {
			// Create a slice of interface{}'s to represent each column
			values := make([]interface{}, len(columns))
			pointers := make([]interface{}, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}

			// Scan the result into the pointers
			err := rows.Scan(pointers...)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to scan row: %v", err))
				continue
			}

			// Create a map for the current row
			rowData := make(map[string]interface{})
			for i, col := range columns {
				val := values[i]
				b, ok := val.([]byte)
				if ok {
					// Convert []byte to string for better handling in templates
					rowData[col] = string(b)
				} else {
					rowData[col] = val
				}
			}

			results = append(results, rowData)
		}

		// Check for errors during iteration
		if err = rows.Err(); err != nil {
			errs = append(errs, fmt.Errorf("error during rows iteration: %v", err))
		}

		// Add results to template context
		queryData := map[string]interface{}{
			"results": results,
			"count":   len(results),
		}

		// Save current context and create new nested context
		ctx.Data["query"] = map[string]interface{}{
			"query": queryData,
		}

		// Execute the inner content with the new context
		contentErrs := core.Render(ctx, sb, content)
		if errs != nil {
			errs = append(errs, contentErrs...)
		}

		// Remove query results
		ctx.Data["query"] = nil

		return newEndPos, errs
	},
}

// You're awesome now we need to implement a simple but comprehensive tag for reading from a sqlite database safely executes queries.
// let's use https://github.com/tursodatabase/go-libsql

// Existing Tag Examples
// ```
// var ImportTag = TemplateTag{
//     Name: "import",
//     Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (int, []error) {
//         var errs []error
//         path := strings.TrimSpace(tag_contents)

//         // Handle ~/ prefix for adjacent files
//         if strings.HasPrefix(path, "~/") {
//             // Resolve relative to current template
//             baseDir := filepath.Dir(ctx.CurrentTemplatePath)
//             path = filepath.Join(baseDir, strings.TrimPrefix(path, "~/"))
//         }

//         // Determine file type
//         ext := strings.ToLower(filepath.Ext(path))

//         // Read file content
//         content, err := os.ReadFile(path)
//         if err != nil {
//             errs = append(errs, fmt.Errorf("failed to read import file %s: %v", path, err))
//             return pos, errs
//         }

//         // Create appropriate asset based on file type
//         switch ext {
//         case ".css":
//             ctx.AssetRegistry.Add(&structure.Asset{
//                 Type:     structure.CSS,
//                 Content:  string(content),
//                 IsInline: true,
//                 Priority: 100, // Default CSS priority
//             })
//         case ".js":
//             ctx.AssetRegistry.Add(&structure.Asset{
//                 Type:     structure.JS,
//                 Content:  string(content),
//                 IsInline: true,
//                 Priority: 200, // Default JS priority
//             })
//         default:
//             errs = append(errs, fmt.Errorf("unsupported import file type: %s", ext))
//         }

//         return pos, errs
//     },
// }

// var IfTag = TemplateTag{
// 	Name: "if",
// 	Render: func(ctx *structure.RenderCtx, sb *strings.Builder, tag_contents, raw string, pos int) (new_pos int, errs []error) {
// 		endTag := delimWrap(ctx, "endif")
// 		endTagStart, endTagLength := core.SeekIndexAndLength(raw, endTag, pos)
// 		if endTagStart == -1 {
// 			errs = append(errs, fmt.Errorf("could not find end tag for %s", endTag))
// 			return pos, errs
// 		}
// 		content := raw[pos:endTagStart]
// 		newEndPos := endTagStart + endTagLength

// 		value, condition_errors := core.ResolveCondition(ctx, tag_contents)
// 		if len(condition_errors) > 0 {
// 			errs = append(errs, condition_errors...)
// 		}
// 		if value {
// 			sb.WriteString(content)
// 		}
// 		return newEndPos, errs
// 	},
// }

// ```

// // Parse query options
// ```
// options := parseAssetTagOptions(tag_contents)
// query := options["query"]
// dbPath := options["path"]
// ````

// ```
// {% sqlite
//     path="@root/data/myapp.db"
//     query="SELECT * FROM products WHERE stock > 0"
// %}
// ```

// ```
// {% sqlite query="SELECT * FROM users WHERE id = ?" {{ user_id }} %}
// ```

// ```
// {% sqlite query="SELECT COUNT(*) as count FROM notifications WHERE user_id = ?" {{ user.id }} %}
//   {% if .query.results[0].count > 0 %}
//     <div class="badge">{{ .query.results[0].count }}</div>
//   {% endif %}
// {% endsqlite %}
// ```
