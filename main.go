// References/Inspiration:
// - https://github.com/cbroglie/mustache/blob/master/mustache.go

package main

import (
	"fmt"
	"kato-studio/katoengine/lib/engine/static"
	"kato-studio/katoengine/lib/engine/template"
	"kato-studio/katoengine/lib/store"
	"kato-studio/katoengine/lib/utils"
	"log"
	"os"
	"runtime"
	"strings"

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
		}
	}`

	app.Get("/static", func(c *fiber.Ctx) error {
		// render all pages and folders
		static.RenderFolder(pagesDir, gjson.Parse(default_data))

		return c.SendString("Rendered all pages and folders")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// timer start to log processing time
		pageBytes, err := os.ReadFile("./view/pages/+page.kato")
		if err != nil {
			utils.Fatal(fmt.Sprint(err))
		}

		fileData := template.Render(string(pageBytes), gjson.Parse(default_data), components.Store())

		c.Set("Content-Type", "text/html")
		return c.Send([]byte(fileData))
	})

	// redirect and get array of all pages and folders
	app.Get("*", func(c *fiber.Ctx) error {
		// get path
		path := c.Path()

		return c.SendFile("./build/pages/" + path + "index.html")
	})

	if windows {
		log.Fatal(app.Listen("localhost:3000"))
	} else {
		log.Fatal(app.Listen(":3000"))
	}
	utils.ServerPrint("Server started on 🚀 http://localhost:3000")
}
