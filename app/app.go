package app

import (
	"strconv"
	"encoding/json"
	"log"

	"github.com/mteam88/un-abandon/database"

	"github.com/gofiber/fiber/v2" // gofiber import
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html" // mustache template import
)

var App *fiber.App
var DB *database.MemDB

func Setup() {
	DB = database.NewMemDB()
	users, err := json.Marshal([]database.User{})
	if err != nil {
		panic(err)
	}
	DB.Set("users", users)

	repos, err := json.Marshal([]database.Repo{})
	if err != nil {
		panic(err)
	}
	DB.Set("abandoned_repos", repos)

	
	// init fiber app
	App = fiber.New(fiber.Config{
		Views: html.New("./web/views", ".html"), // set the views directory
	})

	// use logger middleware
	App.Use(logger.New())
		
	// define route
	App.Get("/", func(c *fiber.Ctx) error {
		// return index.html
		return c.Render("index", fiber.Map{
			"Header": "Un-Abandon",
		}, "layouts/main")
	})
	ExploreSetup()
	InstallSetup()
	DashboardSetup()
	// serve static files
	App.Static("/", "./web/static/public")
	App.Static("/assets", "./web/static/assets")

}

func Start(port int) {
	// start server
	App.Listen(":" + strconv.Itoa(port))
}
