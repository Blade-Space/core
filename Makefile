dev:
	rm -rf apis/
	rm -rf out
	mkdir apis
	go build main.go
	./main
	go get github.com/gin-gonic/gin
	go run out/app.go