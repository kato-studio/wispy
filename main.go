// References/Inspiration:
// - https://github.com/cbroglie/mustache/blob/master/mustache.go

package main

import (
	"fmt"
	"kato-studio/katoengine/engine"
	"kato-studio/katoengine/utils"
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
	var default_data = `{
		"stars": ["STAR", "STAR-STAR", "STAR-STAR-STAR", "STAR-STAR-STAR-STAR", "STAR-STAR-STAR-STAR-STAR"],
		"page": {
			"title": "Home",
			"url": "/page/url/home",
		},
		"links": [
			{"text": "Home", "url": "/"},
			{"text": "About", "url": "/about"},
			{"text": "Contact", "url": "/contact"},
		],
		"clients": ["client1", "client2", "client3"],
		"data": {
			"is_logged_in":"true",
		},
		"company": {
			"name": "Kato",
			"address": "1234 Kato Lane",
		}
	}`

	var performance_data = map[string][]string{
		"/kato": {},
		"/html": {},
		"/test": {},
		"/slip": {},
	}

	app.Get("pref", func(c *fiber.Ctx) error {
		return c.JSON(performance_data)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// timer start to log processing time
		start := time.Now()

		pageBytes, err := os.ReadFile("./view/pages/+page.hstm")
		utils.Fatal(err)

		var ctx = engine.RenderCTX{
			Json:      gjson.Parse(default_data),
			Snippet:   map[string]string{},
			Variables: map[string]string{},
		}

		page := engine.SlipEngine(string(pageBytes), ctx)

		fmt.Println("Processing time: ", time.Since(start))
		performance_data["/slip"] = append(performance_data["/slip"], fmt.Sprint(time.Since(start)))
		return c.SendString(page)
	})

	app.Get("/plain", func(c *fiber.Ctx) error {
		// timer start to log processing time
		start := time.Now()
		//
		pageBytes, err := os.ReadFile("./view/pages/+page.hstm")
		utils.Fatal(err)
		//
		fmt.Println("Processing time: ", time.Since(start))
		performance_data["/plain"] = append(performance_data["/plain"], fmt.Sprint(time.Since(start)))
		return c.SendString(string(pageBytes))
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
