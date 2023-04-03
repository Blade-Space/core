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
