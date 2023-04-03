// Package handlercore provides functions to generate the main application file
// with the required imports and route configurations for the specified APIs.
package handlercore

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Init generates the main application file with the required imports and
// route configurations for the specified APIs, and a given server port.
//
// Parameters:
//   - apiNames ([]string): A slice of API names to be included in the application.
//   - port (string): The port number for the server to listen on.
func Init(apiNames []string, port string) {
	createAppGoFile(apiNames, port)
}

// createAppGoFile generates the main application file with the required imports
// and route configurations for the specified APIs, and a given server port.
// The generated file will be saved as "out/app.go".
//
// Parameters:
//   - apiNames ([]string): A slice of API names to be included in the application.
//   - port (string): The port number for the server to listen on.
//
// Returns:
//   - error: An error that occurred during the file generation or saving, or nil if successful.
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

	"net"
	"log"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	%s

	internalIP, err := GetLocalIP()

	fmt.Println("Server OS has been successfully launched and is ready to go. 🚀")
	fmt.Println("Run on http://localhost:%s")
	if err == nil {
		fmt.Println("Run on: http://" + internalIP + ":%s")
	}
	if err := r.Run(":%s"); err != nil {
		log.Fatalf(err.Error())
	}
}

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		ip := ipNet.IP
		if ip.IsLoopback() || ip.IsUnspecified() || ip.To4() == nil {
			continue
		}

		return ip.String(), nil
	}

	return "", fmt.Errorf("внутренний IP-адрес не найден")
}
`, strings.Join(apiImports, "\n\t"), strings.Join(apiGroups, "\n\t"), port, port, port)

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
