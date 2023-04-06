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
	"embed"
	"io/fs"
	"errors"
	"io"
	"net/http"
	"strings"
)

//go:embed front-end/dist/*
var frontend embed.FS

type readerAt struct {
	io.Reader
}

func (r *readerAt) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errors.New("readerAt: invalid offset")
	}
	n, err = r.Read(p)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Serve static files
	frontendFS, err := fs.Sub(frontend, "front-end/dist")
	if err != nil {
		panic(err)
	}

	// Custom handler for static files and index.html
	// Serve index.html for the root path
	r.GET("/", func(c *gin.Context) {
		fileServer := http.FileServer(http.FS(frontendFS))
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	// Custom handler for static files
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") {
			c.Status(http.StatusNotFound)
			return
		}

		// Serve static files
		fileServer := http.FileServer(http.FS(frontendFS))
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

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
