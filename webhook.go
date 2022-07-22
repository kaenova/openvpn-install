package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	gosrc "github.com/kaenova/openvpn-install/go_src"
)

func main() {
	if !gosrc.IsRoot() {
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

	app.Get("/", gosrc.RenderIndex)

	app.Get("/list", gosrc.ListClient)

	app.Get("/list/:clientname", gosrc.SendConfig)

	app.Post("/add", gosrc.AddOpenVPNClient)

	app.Listen(":3000")
}
