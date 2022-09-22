package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/omarahm3/compat/compat"
)

const (
	COMPAT_FILE = "compat"
)

func main() {
	path, err := os.Getwd()
	must(err)

	filePath, err := getCompatFilePath(path)
	must(err)

	content := readFile(*filePath)

	err = compat.Run(content)
	must(err)
}

func getCompatFilePath(cwd string) (*string, error) {
	ymlFile := fmt.Sprintf("%s/%s.yml", cwd, COMPAT_FILE)
	yamlFile := fmt.Sprintf("%s/%s.yaml", cwd, COMPAT_FILE)

	if _, err := os.Stat(ymlFile); !os.IsNotExist(err) {
		return &ymlFile, nil
	}

	if _, err := os.Stat(yamlFile); !os.IsNotExist(err) {
		return &yamlFile, nil
	}

	return nil, errors.New("Compat file was not found. Please consider creating one.")
}

func readFile(file string) []byte {
	content, err := ioutil.ReadFile(file)
	must(err)

	return content
}

func must(err error) {
	if err == nil {
		return
	}

	fmt.Println(err.Error())
	os.Exit(1)
}
