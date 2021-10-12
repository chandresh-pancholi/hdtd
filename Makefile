.PHONY: build-hdtd
build-hdtd:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o hdtd cmd/main.go