package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/list", listClient)

	app.Get("/list/:username", sendConfig)

	app.Post("/add", addOpenVPNClient)

	app.Post("/revoke/:username", revokeClinet)

	app.Listen(":3000")
}

func addOpenVPNClient(c *fiber.Ctx) error {
	return nil
}

func sendConfig(c *fiber.Ctx) error {
	return nil
}

func listClient(c *fiber.Ctx) error {
	return nil
}

func revokeClinet(c *fiber.Ctx) error {
	return nil
}
