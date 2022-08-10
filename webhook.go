package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/kaenova/openvpn-install/pkg"
)

func main() {
	if !pkg.IsRoot() {
		log.Fatal("You need to run this program on root")
	}

	templateEngine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: templateEngine,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://ops.kerjago.id, localhost:3000, http://127.0.0.1:3000, http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Static("/static", "./config", fiber.Static{
		Download: true,
		Browse:   false,
	})

	app.Get("/", pkg.RenderIndex)

	app.Get("/list", pkg.ListClient)

	app.Post("/add", pkg.AddOpenVPNClient)

	app.Post("/revoke", pkg.RevokeClinet)

	app.Get("/download", pkg.SendConfig)

	app.Listen(":3000")
}
