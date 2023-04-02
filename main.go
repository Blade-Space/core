package main

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

// ! Скачивание Репозиториев
func downloadGitRepositories(urls []string) {
	for _, url := range urls {
		downloadGitRepository(url)
	}
}

func downloadGitRepository(url string) {
	repoName := extractRepoName(url)

	_, err := exec.Command("git", "clone", url, "apis/"+repoName).Output()
	if err != nil {
		log.Fatalf("Ошибка при клонировании репозитория %s: %v", url, err)
	}
	fmt.Printf("Репозиторий %s успешно клонирован.\n", url)
}

func extractRepoName(url string) string {
	parts := strings.Split(url, "/")
	name := strings.TrimSuffix(parts[len(parts)-1], ".git")
	return name
}

// ! Скачивание репозиториев оконченно (END)

// ! Проходимся по папкам api и считываем api.yml
type APIInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	API     string `yaml:"api"`
}

func readAPIInfoFromRepositories(repoPaths []string) {
	for _, repoPath := range repoPaths {
		apiInfo, err := readAPIInfo(filepath.Join(repoPath, "api.yml"))
		os.Remove(filepath.Join(repoPath, "go.mod"))
		os.Remove(filepath.Join(repoPath, "go.sum"))

		if err != nil {
			log.Printf("Ошибка при чтении файла api.yml из %s: %v", repoPath, err)
			continue
		}
		fmt.Printf("Информация об API в %s: %+v\n", repoPath, apiInfo)
	}
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

// ! Проходимся по папкам api и считываем api.yml (END)

// ! Отвратительная реализация хорошей кодогенерации
func createAppGoFile(apiNames []string) error {
	apiImports := make([]string, len(apiNames))
	apiGroups := make([]string, len(apiNames))

	for i, apiName := range apiNames {
		apiImports[i] = fmt.Sprintf(`%s "simplified-prototype-api-collector/apis/api-%s/routes"`, apiName, apiName)
		apiGroups[i] = fmt.Sprintf(`
	api%s := r.Group("/api/%s")
	%s.RegisterRoutes(api%s)`, apiName, apiName, apiName, apiName)
	}

	content := fmt.Sprintf(`package main

import (
	%s

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// API (START)
	%s
	// API (END)

	r.Run(":3000")
}
`, strings.Join(apiImports, "\n\t"), strings.Join(apiGroups, "\n\t"))

	outDir := "out"
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		log.Println(err)
		return err
	}

	filePath := outDir + "/app.go"
	err = ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// ! Отвратительная реализация хорошей кодогенерации (END)

func main() {
	gitUrls := []string{
		"https://github.com/Blade-Space/api-wwf",
	}

	downloadGitRepositories(gitUrls)

	repoNames := []string{
		"./apis/api-wwf",
	}

	readAPIInfoFromRepositories(repoNames)

	apiNames := []string{"wwf"}
	err := createAppGoFile(apiNames)
	if err != nil {
		fmt.Printf("Ошибка при создании файла app.go: %v\n", err)
	}
}
