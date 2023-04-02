package handlercore

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func Init(apiNames []string, port string) {
	createAppGoFile(apiNames, port)
}

func createAppGoFile(apiNames []string, port string) error {
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

	r.Run(":%s")
}
`, strings.Join(apiImports, "\n\t"), strings.Join(apiGroups, "\n\t"), port)

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
