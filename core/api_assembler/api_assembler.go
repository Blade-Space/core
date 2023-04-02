package apiassembler

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type APIInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	API     string `yaml:"api"`
}

func Init(urls []string) []APIInfo {
	fmt.Println("Starting downloads repos")
	var repos []string

	for _, url := range urls {
		repo := downloadGitRepository(url)
		repos = append(repos, repo)
	}

	reposContent := readAPIInfoFromRepositories(repos)

	return reposContent
}

func downloadGitRepository(url string) string {
	repoName := extractRepoName(url)

	_, err := exec.Command("git", "clone", url, "apis/"+repoName).Output()
	if err != nil {
		log.Fatalf("Ошибка при клонировании репозитория %s: %v", url, err)
	}

	return "apis/" + repoName
}

func extractRepoName(url string) string {
	parts := strings.Split(url, "/")
	name := strings.TrimSuffix(parts[len(parts)-1], ".git")
	return name
}

func readAPIInfoFromRepositories(repoPaths []string) []APIInfo {
	var apiContent []APIInfo

	for _, repoPath := range repoPaths {
		apiInfo, err := readAPIInfo(filepath.Join(repoPath, "api.yml"))
		apiContent = append(apiContent, apiInfo)

		os.Remove(filepath.Join(repoPath, "go.mod"))
		os.Remove(filepath.Join(repoPath, "go.sum"))
		os.Remove(filepath.Join(repoPath, "main.go"))

		if err != nil {
			log.Printf("Ошибка при чтении файла api.yml из %s: %v", repoPath, err)
			continue
		}
		fmt.Printf("API в %s: %+v\n", repoPath, apiInfo.Name)
	}

	return apiContent
}

func readAPIInfo(filepath string) (APIInfo, error) {
	var apiInfo APIInfo

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return apiInfo, err
	}

	err = yaml.Unmarshal(data, &apiInfo)
	if err != nil {
		return apiInfo, err
	}

	return apiInfo, nil
}
