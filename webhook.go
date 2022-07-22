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

	app.Get("/list/:clientname", SendConfig)

	app.Post("/add", AddOpenVPNClient)

	app.Listen(":3000")
}

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
	users, err := GetUserActive()
	if err != nil {
		c.Status(500)
		return c.SendString("Tidak dapat mengambil config user")
	}
	return c.JSON(users)
}

// func revokeClinet(c *fiber.Ctx) error {
// 	return nil
// }

func IsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

func GetUserOpenVPNFile() ([]string, error) {
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

func GetUserActive() ([]string, error) {
	usersOpenvpn, err := GetUserOpenVPNFile()
	if err != nil {
		return nil, err
	}
	usersDir, err := ListDir("./config")
	if err != nil {
		return nil, err
	}

	var finalUser []string
	for i := 0; i < len(usersOpenvpn); i++ {
		if StringInSlice(usersOpenvpn[i]+".ovpn", usersDir) {
			finalUser = append(finalUser, usersOpenvpn[i])
		}
	}
	return finalUser, nil
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ListDir(dir string) ([]string, error) {
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
