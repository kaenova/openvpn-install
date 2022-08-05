package main

import (
	"os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

/*
=== Handler ===
*/

func RenderIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

type CreateOpenVPNClientReq struct {
	ClientName string `json:"client" xml:"client" form:"client"`
}

func AddOpenVPNClient(c *fiber.Ctx) error {
	client := new(CreateOpenVPNClientReq)

	users, err := GetUserActive()
	if err != nil {
		c.Status(500)
		return c.SendString("Tidak dapat mengambil config user yang tersedia")
	}

	if err := c.BodyParser(client); err != nil {
		c.Status(400)
		return c.SendString(err.Error())
	}

	if StringInSlice(client.ClientName, users) {
		c.Status(400)
		return c.SendString("Client name sudah digunakan")
	}

	os.Setenv("MENU_OPTION", "1")
	os.Setenv("CLIENT", client.ClientName)
	os.Setenv("PASS", "1")
	cmd := exec.Command("./openvpn-install.sh")
	cmd.Start()
	err = cmd.Wait()
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Status(200).SendString("Successfully create " + client.ClientName)
}

func SendConfig(c *fiber.Ctx) error {
	files, err := ListDir("./config")
	if err != nil {
		c.Status(500)
		return c.SendString("Cannot get config file")
	}

	clientName := c.Params("clientname")

	users, err := GetUserOpenVPNFile()
	if err != nil {
		return c.Status(500).SendString("Cannot get current OpenVPN user")
	}

	if !StringInSlice(clientName, users) {
		return c.Status(400).SendString("User not found in OpenVPN")
	}

	if !StringInSlice(clientName+".ovpn", files) {
		c.Status(400)
		return c.SendString("Tidak ditemukan config ovpn yang dicari")
	}

	return c.Redirect("/static/" + clientName + ".ovpn")
}

func ListClient(c *fiber.Ctx) error {
	users, err := GetUserActive()
	if err != nil {
		c.Status(500)
		return c.SendString("Tidak dapat mengambil config user")
	}
	return c.JSON(users)
}

type RevokeOpenVPNClientReq struct {
	ClientName string `json:"client" xml:"client" form:"client"`
}

func RevokeClinet(c *fiber.Ctx) error {
	client := new(RevokeOpenVPNClientReq)

	users, err := GetUserOpenVPNFile()
	if err != nil {
		return c.Status(500).SendString("Cannot get current OpenVPN user")
	}

	if err := c.BodyParser(client); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if !StringInSlice(client.ClientName, users) {
		return c.Status(400).SendString("User not found in OpenVPN")
	}

	os.Setenv("MENU_OPTION", "2")
	os.Setenv("CLIENT", client.ClientName)
	cmd := exec.Command("./openvpn-install.sh")
	cmd.Start()
	err = cmd.Wait()
	if err != nil {
		c.Status(500)
		return c.SendString(err.Error())
	}

	return c.Status(200).SendString("Successfully revoke " + client.ClientName)
}
