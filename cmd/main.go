package main

import (
	assembler "simplified-prototype-api-collector/core/api_assembler"
	filecore "simplified-prototype-api-collector/core/file_core"
	frontend_assembler "simplified-prototype-api-collector/core/frontend_assembler"
	handler "simplified-prototype-api-collector/core/handler_core"
)

func main() {
	// * Инициализация Config
	config := filecore.Init()

	// * Инициализация API Assembler
	repoNames := assembler.Init(config.Repos)

	var names []string

	for _, repo := range repoNames {
		names = append(names, repo.API)
	}

	// * Инициализация Handler
	handler.Init(names, config.Port)

	// * Инициализация Front-end Assembler
	if config.FrontEnd.Include {
		frontend_assembler.Init(config.FrontEnd.Repo)
	}
}
