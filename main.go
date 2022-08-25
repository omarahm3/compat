package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/omarahm3/compat/compat"
)

func main() {
	flag.Parse()

	filePath, err := getCompatFilePath()
	must(err)

	content := readFile(*filePath)

	err = compat.Run(content)
	must(err)
}

func getCompatFilePath() (*string, error) {
	path, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	var f string

	ymlPath, err := checkCompatFile(fmt.Sprintf("%s/compat.yml", path))

	if errors.Is(err, os.ErrNotExist) {
		yamlPath, err := checkCompatFile(fmt.Sprintf("%s/compat.yaml", path))
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("compat file does not exist")
		}

		f = *yamlPath
	} else {
		f = *ymlPath
	}

	return &f, nil
}

func checkCompatFile(path string) (*string, error) {
	_, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	return &path, nil
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
