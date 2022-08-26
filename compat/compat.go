package compat

import (
	"fmt"
	"os"
	"strings"

	yaml "github.com/goccy/go-yaml"
)

const (
	base_prefix   = "base_"
	extends_token = "inherit"
)

func Run(content []byte) error {
	var yamlFile yaml.MapSlice
	err := yaml.Unmarshal(content, &yamlFile)
	must(err)

	for _, v := range yamlFile {
		if v.Key.(string) == "services" {
			processServices(&v)
		}
	}

	write(&yamlFile)

	return nil
}

func write(data *yaml.MapSlice) {
	var buf strings.Builder
	e := yaml.NewEncoder(&buf, yaml.UseSingleQuote(true), yaml.Indent(2))

	err := e.Encode(data)
	must(err)

	path, err := os.Getwd()
	must(err)

	writeFile(buf, fmt.Sprintf("%s/docker-compose.yaml", path))
}

func writeFile(data strings.Builder, composePath string) {
	f, err := os.Create(composePath)
	must(err)
	defer f.Close()

	_, err = f.WriteString(data.String())
	must(err)
}

func processServices(services *yaml.MapItem) {
	baseServices := make(map[string]interface{})
	svcs := services.Value.(map[string]interface{})

	// Build base services map
	for key, service := range svcs {
		if strings.Contains(key, "base_") {
			baseServices[key] = service
			delete(svcs, key)
		}
	}

	// Process services
	for _, service := range svcs {
		svc := service.(map[string]interface{})
		for k, baseServiceName := range svc {
			if k == extends_token {
				delete(svc, k)
				baseService := baseServices[baseServiceName.(string)].(map[string]interface{})

				for baseServiceKey, baseServiceValue := range baseService {
					svc[baseServiceKey] = baseServiceValue
				}
			}
		}
	}
}

func must(err error) {
	if err == nil {
		return
	}

	fmt.Println(err.Error())
	os.Exit(1)
}
