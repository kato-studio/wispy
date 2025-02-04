package engine

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

// ---====----
// Engine Render/Template functions
// ---====----
func (e *EngineCtx) RenderRoute(w io.Writer, host string, data map[string]interface{}, c echo.Context) error {
	name := strings.TrimSuffix(c.Request().URL.Path, "/")
	templateName := host + filepath.FromSlash(name)

	// Get the specific template
	site, siteMapOk := e.SiteMap[host]
	route, routeExist := site.Routes[templateName]
	if !siteMapOk || !routeExist {
		fmt.Println("Site not found:", host, "(siteFound?:", siteMapOk, "routeExist:", routeExist, ")")
		fmt.Println("~looking for", templateName, "current routes... [ ]")
		return echo.NewHTTPError(404, "404 not found,", c.Request().URL.Path)
	}

	// Create a new rendering context
	ctx := NewRenderCtx(e, data, site)

	// ✅ Ensure FuncMap is applied before parsing the template
	// base, baseErr := template.New("base").Funcs(GetDefaultFuncs(&ctx)).Parse(route.Template)
	// if baseErr != nil {
	// 	fmt.Println("Error parsing template:", baseErr)
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "Error parsing template:", route.Path, "\n\n", baseErr)
	// }

	// // Render the template into a buffer
	// renderedPage := bytes.NewBuffer([]byte{})
	// err := base.Execute(renderedPage, data)
	// if err != nil {
	// 	fmt.Println("Error executing template:", err)
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "Error executing template: \n"+fmt.Sprint(err), route.Path, "\n\n", err)
	// }

	renderedPage := bytes.NewBuffer([]byte{})
	// Store the rendered page content
	data["PageContents"] = renderedPage.String()
	data["PageHeadContents"] = e.CreateHeadContent(route)

	// 🚀 Final output
	fmt.Fprintf(w, `<!DOCTYPE html><html lang="%s">%s<body>%s</body></html>`, ctx.Lang, data["PageHeadContents"], data["PageContents"])
	return nil
}

func (e *EngineCtx) CreateHeadContent(route PageRoutes) string {
	mt := route.MetaTags

	ogTitle := mt.OgTitle
	if ogTitle == "" {
		ogTitle = mt.Title
	}
	ogDescription := mt.OgDescription
	if ogDescription == "" {
		ogDescription = mt.OgDescription
	}

	templateString := fmt.Sprintf("<title>%s</title>", mt.Title) +
		fmt.Sprintf("<meta name=\"%s\" contents=\"%s\">", "title", mt.Title) +
		fmt.Sprintf("<meta name=\"%s\" contents=\"%s\">", "description", mt.Description) +
		fmt.Sprintf("<meta name=\"%s\" contents=\"%s\">", "og:title", ogTitle) +
		fmt.Sprintf("<meta name=\"%s\" contents=\"%s\">", "og:description", ogDescription) +
		fmt.Sprintf("<meta name=\"%s\" contents=\"%s\">", "og:type", mt.OgType) +
		fmt.Sprintf("<meta name=\"%s\" contents=\"%s\">", "og:url", mt.OgUrl)

	for name, value := range mt.OtherTags {
		templateString += fmt.Sprintf("<meta name=\"%s\" contents=\"%s\">", name, value)
	}

	return templateString

}
