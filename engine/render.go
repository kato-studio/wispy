package engine

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"maps"
	"net/http"
	"os"
	"strings"
	textDanger "text/template"

	"github.com/labstack/echo/v4"
)

// ---====----
// Engine Render/Template functions
// ---====----
func (e *EngineCtx) RenderRoute(w io.Writer, host string, data map[string]interface{}, c echo.Context) error {
	name := strings.TrimSuffix(c.Request().URL.Path, "/")
	templateName := host + name
	// Get the specific template
	site, siteMapOk := e.SiteMap[host]
	route, routeExist := site.Routes[templateName]
	if !siteMapOk || !routeExist {
		// TODO: error reporting
		fmt.Println("~looking for", templateName)
		e.Log.Error("Site not found: ", host)
		fmt.Print(maps.Keys(site.Routes))
		return echo.NewHTTPError(404, "404 not found,", c.Request().URL.Path)
	}
	//
	base, baseErr := template.New("base").Funcs(GetDefaultFuncs(e)).Parse(route.Template)
	if baseErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Create page template", route.Path, "\n\n", baseErr)
	}

	renderedPage := bytes.NewBuffer([]byte{})
	base.Execute(renderedPage, data)
	data["PageContents"] = string(renderedPage.String())

	// Todo: only handle this if needed otherwise if flag to return only page contents
	data["PageHeadContents"] = e.CreateHeadContent(route)

	//
	layoutBytes, err := os.ReadFile(route.Layout)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could read layout file", route.Layout, "\n\n", err)
	}

	layoutTemplate, err := textDanger.New(templateName + ":layout").Parse(string(layoutBytes))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could render with layout", templateName, err)
	}

	// 🚀 🚀
	return layoutTemplate.Execute(w, data)
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
