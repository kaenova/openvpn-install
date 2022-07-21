package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if !isRoot() {
		log.Fatal("You need to run this program on root")
	}

	app := fiber.New()

	app.Static("/static", "./config", fiber.Static{
		Download: true,
		Browse:   false,
	})

	app.Get("/list", listClient)

	app.Get("/list/:clientname", sendConfig)

	app.Post("/add", addOpenVPNClient)

	// app.Post("/revoke/:username", revokeClinet)

	app.Listen(":3000")
}

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

func getUser() ([]string, error) {
	readFile := "/etc/openvpn/easy-rsa/pki/index.txt"
	file, err := os.Open(readFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		re := regexp.MustCompile(`/CN=([^\s]+)`)
		result := re.FindStringSubmatch(line)

		if len(result) == 2 {
			users = append(users, result[1])
		}

	}
	return users, scanner.Err()
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func listDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil
	}

	var allFiles []string
	for _, f := range files {
		allFiles = append(allFiles, f.Name())
	}
	return allFiles, nil
}

/*
=== Handler ===
*/

type CreateOpenVPNClientReq struct {
	ClientName string  `json:"client" xml:"client" form:"client"`
	Password   *string `json:"password" xml:"password" form:"password"`
}

func addOpenVPNClient(c *fiber.Ctx) error {
	client := new(CreateOpenVPNClientReq)

	users, err := getUser()
	if err != nil {
		c.Status(500)
		return c.SendString("Tidak dapat mengambil config user yang tersedia")
	}

	if err := c.BodyParser(client); err != nil {
		c.Status(400)
		return c.SendString(err.Error())
	}

	if stringInSlice(client.ClientName, users) {
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

func sendConfig(c *fiber.Ctx) error {
	files, err := listDir("./config")
	if err != nil {
		c.Status(500)
		return c.SendString("Tidak dapat mengambil data config yang tersedia")
	}

	clientName := c.Params("clientname")
	if !stringInSlice(clientName+".ovpn", files) {
		c.Status(400)
		return c.SendString("Tidak ditemukan config ovpn yang dicari")
	}

	return c.Redirect("/static/" + clientName + ".ovpn")
}

func listClient(c *fiber.Ctx) error {
	users, err := getUser()
	if err != nil {
		c.Status(500)
		return c.SendString("Tidak dapat mengambil config user")
	}
	return c.JSON(users)
}

// func revokeClinet(c *fiber.Ctx) error {
// 	return nil
// }
