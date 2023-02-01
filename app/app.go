package app

import (
	"encoding/json"
	"strconv"

	"github.com/dgraph-io/badger/v3"
	"github.com/mteam88/un-abandon/database"

	"github.com/gofiber/fiber/v2" // gofiber import
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html" // mustache template import
)

var App *fiber.App
var DB *badger.DB

func Setup() {
	DB = database.NewBadgerDB()
	users, err := json.Marshal([]database.User{})
	if err != nil {
		panic(err)
	}
	err = DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("users"), users)
		return err
	  })
	if err != nil {
		panic(err)
	}

	repos, err := json.Marshal([]database.Repo{})
	if err != nil {
		panic(err)
	}
	err = DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("abandoned_repos"), repos)
		return err
	  })
	if err != nil {
		panic(err)
	}

	
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

	App.Use(func(c *fiber.Ctx) error {
		return c.Status(404).Render("404", fiber.Map{
			"Header": "404",
		}, "layouts/main")
	})

}

func Start(port int) {
	// start server
	App.Listen(":" + strconv.Itoa(port))
}
