// References/Inspiration:
// - https://github.com/cbroglie/mustache/blob/master/mustache.go

package main

import (
	"fmt"
	templ "html/template"
	"kato-studio/katoengine/lib/engine/static"
	"kato-studio/katoengine/lib/engine/template"
	"kato-studio/katoengine/lib/store"
	"kato-studio/katoengine/lib/utils"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"

)

func main() {
	// Load .env file
	envRrr := godotenv.Load("./.env")
	if envRrr != nil {
		utils.Fatal("Error loading .env file")
	}

	// check if os is windows
	windows := false
	if runtime.GOOS == "windows" {
		windows = true
	}
	if windows {
		utils.Info("Windows OS detected")
	} else {
		utils.Info("Non-Windows OS detected")
	}

	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		utils.Info("Debug mode enabled")
	} else if strings.ToLower(os.Getenv("DEBUG")) == "false" && windows {
		utils.Info("Debug mode was not set to false & windows has been detected")
		utils.Info("Enabling debug mode....")
	}

	// ------------------------------
	// Begin setup logic for server
	// ------------------------------

	// Default config
	app := fiber.New()

	// render pages based on folder structure in /pages
	const pagesDir = "./view/pages"

	// components store
	components := store.GlobalByteMap()

	// load all components
	const componentsDir = "./view/components/"
	utils.Debug(fmt.Sprintf("Loading components from %s", componentsDir))
	static.LoadAllComponents(componentsDir)
	
	var default_data = `{
		"stars": ["STAR", "STAR-STAR", "STAR-STAR-STAR", "STAR-STAR-STAR-STAR", "STAR-STAR-STAR-STAR-STAR"],
		"page": {
			"title": "Home",
			"url": "/page/url/home",
		},
		"links":["link1", "link2", "link3"],
		"clients": ["client1", "client2", "client3"],
		"data": {
			"is_logged_in":"true",
		},
		"company": {
			"name": "Kato",
			"address": "1234 Kato Lane",
		}
		"Footer": {
			"links": ["link1", "link2", "link3"],
			"company": {
				"name": "Kato",
			}
		},
		"Header": {
			"links": ["link1", "link2", "link3"],
		}
	}`

	var preference_data = map[string][]string{
		"/kato": {},
		"/html": {},
		"/test": {},
	}

	app.Get("pref", func(c *fiber.Ctx) error {
		return c.JSON(preference_data)
	})


	app.Get("/test", func(c *fiber.Ctx) error {
		// timer start to log processing time
		start := time.Now()

		json_data := gjson.Parse(default_data)
		page_bytes, err := os.ReadFile("./view/pages/+page.kato")
		utils.Fatal(err)

		rendered_page := template.SlipEngine(page_bytes, json_data)
		

		c.Set("Content-Type", "text/html")
		fmt.Println("Processing time: ", time.Since(start))
		preference_data["/test"] = append(preference_data["/test"], fmt.Sprint(time.Since(start)))
		return c.SendString(rendered_page)
	})

	


	app.Get("/static", func(c *fiber.Ctx) error {
		// render all pages and folders
		static.RenderFolder(pagesDir, gjson.Parse(default_data))


		return c.SendString("Rendered all pages and folders")
	})

	app.Get("/kato", func(c *fiber.Ctx) error {
		// timer start to log processing time
		start := time.Now()

		pageBytes, err := os.ReadFile("./view/pages/+page.kato")
		if err != nil {
			utils.Fatal(fmt.Sprint(err))
		}

		fileData := template.Render(string(pageBytes), gjson.Parse(default_data), components.Store())

		c.Set("Content-Type", "text/html")
		fmt.Println("Processing time: ", time.Since(start))
		preference_data["/kato"] = append(preference_data["/kato"], fmt.Sprint(time.Since(start)))
		return c.Send([]byte(fileData))
	})

	app.Get("/html", func(c *fiber.Ctx) error {
		// timer start to log processing time
		start := time.Now()
		out := new(strings.Builder)

		page_bytes, _ := os.ReadFile("./templates/pages/page.html")
		page_html := string(page_bytes)
		json_data := gjson.Parse(default_data)

		page_html = template.LoadTemplateComponents(page_html, []string{
			"Header.html",
			"Footer.html",
		})

		utils.Print("html")
		utils.Print(page_html)

		page_template, err := templ.New("page").Parse(page_html)
		if err != nil {
			utils.Fatal(fmt.Sprint(err))
		}
		utils.Print(json_data.Value())
		err = page_template.ExecuteTemplate(out, "page", json_data.Value())
		if err != nil { utils.Fatal(fmt.Sprint(err)) }

		c.Set("Content-Type", "text/html")
		fmt.Println("Processing time: ", time.Since(start))
		preference_data["/html"] = append(preference_data["/html"], fmt.Sprint(time.Since(start)))
		return c.SendString(out.String())
	})

	// redirect and get array of all pages and folders
	app.Get("*", func(c *fiber.Ctx) error {
		// get path
		path := c.Path()
		// USE SSR? check if there is and index.server.kato file
		if _, err := os.Stat("./build/pages/" + path + "server.kato"); err == nil {
			pageBytes, err := os.ReadFile("./build/pages/" + path + "server.kato")

			if err != nil {
				utils.Fatal(fmt.Sprint(err))
			}

			fileData := template.Render(string(pageBytes), gjson.Parse(default_data), components.Store())

			c.Set("Content-Type", "text/html")
			return c.Send([]byte(fileData))
		}
		// server static files
		return c.SendFile("./build/pages/" + path + "index.html")
	})

	// this windows check is to prevent the server from failing to bind to the port on windows
	if windows {
		log.Fatal(app.Listen("localhost:3000"))
		utils.ServerPrint("Server started on 🚀 http://localhost:3000")
	} else {
		log.Fatal(app.Listen(":3000"))
		utils.ServerPrint("Server started on 🚀 http://0.0.0.0:3000")
	}
}
