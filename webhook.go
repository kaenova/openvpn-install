package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	if !IsRoot() {
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

	app.Get("/", RenderIndex)

	app.Get("/list", ListClient)

	app.Post("/add", AddOpenVPNClient)

	app.Post("/invoke", RevokeClinet)

	app.Get("/download", SendConfig)

	app.Listen(":3000")
}
