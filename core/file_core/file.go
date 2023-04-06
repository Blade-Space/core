package filecore

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type FrontEnd struct {
	Include bool   `yaml:"include"`
	Repo    string `yaml:"repo"`
}

type ReposConfig struct {
	Name     string   `yaml:"name"`
	Versin   string   `yaml:"versin"`
	Port     string   `yaml:"port"`
	Repos    []string `yaml:"repos"`
	FrontEnd FrontEnd `yaml:"front-end"`
}

func readReposFromFile(filename string) (ReposConfig, error) {
	var config ReposConfig

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func Init() ReposConfig {
	filename := "config.yaml"

	config, err := readReposFromFile(filename)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла %s: %v", filename, err)
	}

	return config
}
