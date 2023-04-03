package apiassembler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

type APIInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	API     string `yaml:"api"`
}

type APIYaml struct {
	Name    string   `yaml:"name"`
	API     string   `yaml:"api"`
	Link    string   `yaml:"link"`
	Version string   `yaml:"version"`
	Date    string   `yaml:"date"`
	Type    string   `yaml:"type"`
	Authors []string `yaml:"authors"`
}

func FetchAndParseAPIYaml(apiName string) (APIYaml, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/Blade-Space/core/main/packages/%s.yaml", apiName)

	resp, err := http.Get(url)
	if err != nil {
		return APIYaml{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return APIYaml{}, fmt.Errorf("ошибка получения YAML-файла: %s", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIYaml{}, err
	}

	var apiYaml APIYaml
	err = yaml.Unmarshal(data, &apiYaml)
	if err != nil {
		return APIYaml{}, err
	}

	return apiYaml, nil
}

func Init(urls []string) []APIInfo {
	fmt.Println("")
	fmt.Println("Starting downloads repos ⬇️")
	var repos []string

	var wg sync.WaitGroup
	reposChan := make(chan string, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			repo := downloadGitRepository(u)
			reposChan <- repo
		}(url)
	}

	wg.Wait()
	close(reposChan)

	for repo := range reposChan {
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

		apiYaml, err := FetchAndParseAPIYaml(apiInfo.API)
		if err != nil {
			fmt.Println("Ошибка:", err)
			os.Remove(repoPath)
			continue
		}
		fmt.Printf("✅ API YAML: %+v (%+v-%+v) \n", apiYaml.Name, apiYaml.API, apiYaml.Version)

		apiContent = append(apiContent, apiInfo)

		os.Remove(filepath.Join(repoPath, "go.mod"))
		os.Remove(filepath.Join(repoPath, "go.sum"))
		os.Remove(filepath.Join(repoPath, "main.go"))

		if err != nil {
			log.Printf("Ошибка при чтении файла api.yml из %s: %v", repoPath, err)
			continue
		}
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
