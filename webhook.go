package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println(isRoot())
	app := fiber.New()

	app.Get("/list", listClient)

	app.Get("/list/:username", sendConfig)

	app.Post("/add", addOpenVPNClient)

	app.Post("/revoke/:username", revokeClinet)

	app.Listen(":3000")
}

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
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
