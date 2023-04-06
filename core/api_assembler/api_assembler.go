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

// FetchAndParseAPIYaml downloads a YAML file for the specified apiName from a remote
// repository and parses it into an APIYaml structure. If an error occurs while
// downloading or parsing the file, the function returns an error.
//
// Parameters:
// 	- apiName (string): The name of the API for which to download the YAML file.
//
// Returns:
// 	- APIYaml: A structure containing data from the downloaded YAML file.
// 	- error: An error that occurred during the download or parsing of the file, or nil if successful.
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

// Init downloads repositories from the list of URLs and extracts API information from them.
// The function downloads repositories in parallel and returns a slice with API information.
//
// Parameters:
//   - urls ([]string): Slice of repository URLs to download.
//
// Returns:
//   - []APIInfo: A slice of API information extracted from the downloaded repositories.
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

// downloadGitRepository clones a remote git repository using the provided URL.
// The function returns the path to the cloned local repository.
//
// Parameters:
//   - url (string): The URL of the git repository to clone.
//
// Returns:
//   - string: The path to the cloned local repository.
func downloadGitRepository(url string) string {
	repoName := extractRepoName(url)

	_, err := exec.Command("git", "clone", url, "apis/"+repoName).Output()
	if err != nil {
		log.Fatalf("Ошибка при клонировании репозитория %s: %v", url, err)
	}

	return "apis/" + repoName
}

// extractRepoName extracts the repository name from the URL.
//
// Parameters:
//   - url (string): The URL of the git repository.
//
// Returns:
//   - string: The repository name.
func extractRepoName(url string) string {
	parts := strings.Split(url, "/")
	name := strings.TrimSuffix(parts[len(parts)-1], ".git")
	return name
}

// readAPIInfoFromRepositories reads API information from the repositories at the given paths and returns
// a slice of APIInfo structures. It also removes unnecessary files after processing.
//
// Parameters:
//   - repoPaths ([]string): A slice of paths to the repositories to read the API information from.
//
// Returns:
//   - []APIInfo: A slice of APIInfo structures containing the API information from the repositories.
func readAPIInfoFromRepositories(repoPaths []string) []APIInfo {
	var apiContent []APIInfo

	for _, repoPath := range repoPaths {
		apiInfo, err := readAPIInfo(filepath.Join(repoPath, "api.yml"))
		if err != nil {
			fmt.Println("Ошибка:", err)
			os.Remove(repoPath)
			continue
		}

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

// readAPIInfo reads API information from a given file and returns an APIInfo structure.
//
// Parameters:
//   - filepath (string): The path to the file containing the API information.
//
// Returns:
//   - APIInfo: A structure containing the API information from the file.
//   - error: An error that occurred during the reading or parsing of the file, or nil if successful.
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
