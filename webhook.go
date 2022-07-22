package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	if !isRoot() {
		log.Fatal("You need to run this program on root")
	}

	templateEngine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: templateEngine,
	})

	app.Static("/static", "./config", fiber.Static{
		Download: true,
		Browse:   false,
	})

	app.Get("/", renderIndex)

	app.Get("/list", listClient)

	app.Get("/list/:clientname", sendConfig)

	app.Post("/add", addOpenVPNClient)

	app.Listen(":3000")
}
