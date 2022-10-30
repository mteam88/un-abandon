package main

import (
	"github.com/mteam88/un-abandon/app"
)

func main() {
	app.Setup()
	app.Start(3000)
}
