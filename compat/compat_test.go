package compat

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestProcessServices(t *testing.T) {
	strContent := `version: '3.2'
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
    image: mongo`
  expected := `version: '3.2'
services:
  mongo:
    container_name: database
    deploy:
      resources:
        limits:
          cpus: '0.1'
          memory: 128M
        reservations:
          memory: 128M
    image: mongo`
	content := getYamlContent(strContent)

	for _, v := range content {
		if v.Key.(string) == "services" {
			processServices(&v)
		}
	}

  actual := yamlToString(&content)

	if !reflect.DeepEqual(strings.TrimSpace(expected), strings.TrimSpace(actual.String())) {
    t.Errorf("expected:\n%q\nactual:\n%q", expected, actual.String())
  }
}

func TestWriteDockerComposeFile(t *testing.T) {
	strContent := `version: '3.2'
services:
  mongo:
    container_name: database
    deploy:
      resources:
        limits:
          cpus: '0.1'
          memory: 128M
        reservations:
          memory: 128M
    image: mongo`
	content := getYamlContent(strContent)

	write(&content)

	path, _ := os.Getwd()
	filePath := fmt.Sprintf("%s/docker-compose.yaml", path)
	c, _ := ioutil.ReadFile(filePath)

	if !reflect.DeepEqual(strings.TrimSpace(strContent), strings.TrimSpace(string(c))) {
		t.Error("content doesn't match")
	}

	if err := removeFile(filePath); err != nil {
		t.Error("couldn't remove docker-compose.yaml file")
	}
}

func getYamlContent(c string) yaml.MapSlice {
	var content yaml.MapSlice

	yaml.Unmarshal([]byte(c), &content)

	return content
}

func removeFile(path string) error {
	return os.Remove(path)
}
