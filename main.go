// References/Inspiration:
// - https://github.com/cbroglie/mustache/blob/master/mustache.go

package main

import (
	"fmt"
	"kato-studio/katoengine/lib/engine"
	"kato-studio/katoengine/lib/engine/static"
	"kato-studio/katoengine/lib/store"
	"kato-studio/katoengine/lib/utils"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
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
	const componentsDir = "./view/components"
	static.LoadAllComponents(componentsDir)

	app.Get("static", func(c *fiber.Ctx) error {
		// render all pages and folders
		var empty interface{}
		static.RenderFolder(pagesDir, empty)

		return c.SendString("Rendered all pages and folders")

	})

	app.Get("/", func(c *fiber.Ctx) error {
		// timer start to log processing time
		var empty interface{}

		pageBytes, err := os.ReadFile("./view/pages/+page.kato")
		if err != nil {
			utils.Fatal(fmt.Sprint(err))
		}

		fileData := engine.Render(pageBytes, empty, components.Store())

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
