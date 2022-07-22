package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"regexp"
)

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
