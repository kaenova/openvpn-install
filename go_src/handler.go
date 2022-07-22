package gosrc

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
	ClientName string  `json:"client" xml:"client" form:"client"`
	Password   *string `json:"password" xml:"password" form:"password"`
}

func AddOpenVPNClient(c *fiber.Ctx) error {
	client := new(CreateOpenVPNClientReq)

	users, err := GetUser()
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
	os.Setenv("PASS", "")
	if client.Password != nil {
		os.Setenv("PASS", *client.Password)
	}

	cmd := exec.Command("bash", "./openvpn-install.sh")
	cmd.Start()
	err = cmd.Wait()
	if err != nil {
		c.Status(500)
		return c.SendString(err.Error())
	}
	return c.Redirect("/list/" + client.ClientName)
}

func SendConfig(c *fiber.Ctx) error {
	files, err := ListDir("./config")
	if err != nil {
		c.Status(500)
		return c.SendString("Tidak dapat mengambil data config yang tersedia")
	}

	clientName := c.Params("clientname")
	if !StringInSlice(clientName+".ovpn", files) {
		c.Status(400)
		return c.SendString("Tidak ditemukan config ovpn yang dicari")
	}

	return c.Redirect("/static/" + clientName + ".ovpn")
}

func ListClient(c *fiber.Ctx) error {
	users, err := GetUser()
	if err != nil {
		c.Status(500)
		return c.SendString("Tidak dapat mengambil config user")
	}
	return c.JSON(users)
}

// func revokeClinet(c *fiber.Ctx) error {
// 	return nil
// }
