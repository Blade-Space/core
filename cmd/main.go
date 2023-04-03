package main

import (
	assembler "simplified-prototype-api-collector/core/api_assembler"
	filecore "simplified-prototype-api-collector/core/file_core"
	handler "simplified-prototype-api-collector/core/handler_core"
)

func main() {
	// * Initialize Config
	// Read the configuration file that contains the list of repository URLs and the server port.
	config := filecore.Init()

	// * Initialize API Assembler
	// Download the repositories from the list of URLs in the configuration, and extract API information.
	repoNames := assembler.Init(config.Repos)

	var names []string

	// Create a slice of API names to be used in the handler initialization.
	for _, repo := range repoNames {
		names = append(names, repo.API)
	}

	// * Initialize Handler
	// Generate the application's main file (app.go) with the necessary imports and route configurations.
	handler.Init(names, config.Port)
}
