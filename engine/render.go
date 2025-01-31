package engine

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	textDanger "text/template"

	"github.com/labstack/echo/v4"

	"github.com/kato-studio/wispy/engine/templateFuncs"
)

// ---====----
// Engine Render/Template functions
// ---====----
func (e *Engine) RenderRoute(w io.Writer, host string, data map[string]interface{}, c echo.Context) error {
	name := strings.TrimSuffix(c.Request().URL.Path, "/")
	templateName := host + name
	// Get the specific template
	site, siteMapOk := e.SiteMap[host]
	route, routeExist := site.Routes[templateName]
	if !siteMapOk || !routeExist {
		// TODO: error reporting
		e.Log.Error("Site not found: ", host)
		return echo.NewHTTPError(404, "404 not found,", c.Request().URL.Path)
	}

	base := template.New(templateName).Funcs(templateFuncs.GetDefaults())
	base.ParseFiles(route.Template)
	fmt.Println("Components")
	// base.
	// fmt.Println(site.Components.DefinedTemplates())
	// fmt.Println("DefinedTemplates")
	// fmt.Println(base.DefinedTemplates())
	// newBase, errBase := base.AddParseTree("components", site.Components.Tree)
	// fmt.Println(newBase.DefinedTemplates())
	// if errBase != nil {
	// 	fmt.Println(errBase)
	// }

	// fmt.Print("base AFTER ")
	// fmt.Println(base.DefinedTemplates())
	// fmt.Println("---")

	tempWriter := &bytes.Buffer{}

	err := base.ExecuteTemplate(tempWriter, "Page", data)
	if err != nil {
		e.Log.Error("[Error] Could not render ", err)
		fmt.Println("----")
		fmt.Println(err)
		fmt.Println("----")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not render "+templateName, err)
	}

	data["PageHeadContents"] = e.CreateHeadContent(route)
	data["PageContents"] = tempWriter.String()

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

func (e *Engine) CreateHeadContent(route PageRoutes) string {
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
