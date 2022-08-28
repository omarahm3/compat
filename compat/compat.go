package compat

import (
	"fmt"
	"os"
	"reflect"
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

func yamlToString(data *yaml.MapSlice) strings.Builder {
	var buf strings.Builder
	e := yaml.NewEncoder(&buf, yaml.UseSingleQuote(true), yaml.Indent(2))

	err := e.Encode(data)
	must(err)

	return buf
}

func write(data *yaml.MapSlice) {
	buf := yamlToString(data)

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
	baseServices := buildBaseService(services)
	svcs := services.Value.(map[string]interface{})

	// Process services
	for _, s := range svcs {
		service := s.(map[string]interface{})
		for serviceKey, serviceName := range service {
			if serviceKey != extends_token {
				continue
			}

			delete(service, serviceKey)
			inheritedServices := getBaseServicesSlice(serviceName)

			for _, s := range inheritedServices {
				baseService := baseServices[s].(map[string]interface{})

				for baseServiceKey, baseServiceValue := range baseService {
					service[baseServiceKey] = baseServiceValue
				}
			}
		}
	}
}

func buildBaseService(services *yaml.MapItem) map[string]interface{} {
	baseServices := make(map[string]interface{})
	svcs := services.Value.(map[string]interface{})

	// Build base services map
	for key, service := range svcs {
		if strings.Contains(key, "base_") {
			baseServices[key] = service
			delete(svcs, key)
		}
	}

	return baseServices
}

func getBaseServicesSlice(baseServices interface{}) []string {
	t := reflect.TypeOf(baseServices).Kind()
	var bSvcs []string

	if t == reflect.String {
		bSvcs = append(bSvcs, baseServices.(string))
	} else if t == reflect.Slice {
		for _, s := range baseServices.([]interface{}) {
			bSvcs = append(bSvcs, s.(string))
		}
	} else {
		must(fmt.Errorf("unknown type: [%s]", t.String()))
	}

	return bSvcs
}

func must(err error) {
	if err == nil {
		return
	}

	fmt.Println(err.Error())
	os.Exit(1)
}
