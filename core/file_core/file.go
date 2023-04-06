// Package filecore provides functions to read repository configuration from
// a YAML file and return the configuration as a ReposConfig structure.
package filecore

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// readReposFromFile reads the repository configuration from the given file
// and returns a ReposConfig structure. If an error occurs during reading or
// parsing the file, the function returns an error.
//
// Parameters:
//   - filename (string): The name of the file containing the repository configuration.
//
// Returns:
//   - ReposConfig: A structure containing the repository configuration from the file.
//   - error: An error that occurred during the reading or parsing of the file, or nil if successful.
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

// Init reads the repository configuration from a file named "config.yaml" and
// returns a ReposConfig structure. If an error occurs during the reading or
// parsing of the file, the function logs a fatal error and the program will exit.
//
// Returns:
//   - ReposConfig: A structure containing the repository configuration from the file.
func Init() ReposConfig {
	filename := "config.yaml"

	config, err := readReposFromFile(filename)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", filename, err)
	}

	return config
}
