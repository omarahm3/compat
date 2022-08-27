package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

var fileContent = `
version: '3.2'
services:
base_service:
  deploy:
    resources:
      limits:
        cpus: '0.1'
        memory: 128M
      reservations:
        memory: 128M
mongo:
  inherit: base_service
  container_name: database
  image: mongo
`

func InitFile(ext string) error {
	err := ioutil.WriteFile(fmt.Sprintf("%s.%s", COMPAT_FILE, ext), []byte(fileContent), 0644)

	if err != nil {
		return err
	}

	return nil
}

func RemoveFile(ext string) error {
	return os.Remove(fmt.Sprintf("%s.%s", COMPAT_FILE, ext))
}

func TestCompatFileExtensions(t *testing.T) {
	extensions := []string{"yml", "yaml"}

	for _, ext := range extensions {
		if err := InitFile(ext); err != nil {
			t.Error("Couldn't create Compat file")
		}

		cwd, _ := os.Getwd()

		t.Run(fmt.Sprintf("should get correct file path for ext: %s", ext), func(it *testing.T) {
			path, err := getCompatFilePath(cwd)

			if err != nil {
				it.Errorf("error getting compat file: %s", err.Error())
			}

			if !strings.Contains(*path, fmt.Sprintf("%s.%s", COMPAT_FILE, ext)) {
				it.Errorf("got the wrong compat file: %s", *path)
			}
		})

		t.Run(fmt.Sprintf("should read the correct content for ext: %s", ext), func(it *testing.T) {
			content := readFile(fmt.Sprintf("%s/compat.%s", cwd, ext))

			if !reflect.DeepEqual(string(content), fileContent) {
				it.Error("content does not match file content")
			}
		})

		if err := RemoveFile(ext); err != nil {
			t.Errorf("Couldn't remove comapt.%s", ext)
		}
	}
}

func TestCompatFileNotFound(t *testing.T) {
	cwd, _ := os.Getwd()

	t.Run("should return error", func(it *testing.T) {
		_, err := getCompatFilePath(cwd)

		if err == nil {
			it.Errorf("expected error to be returned")
		}
	})
}
