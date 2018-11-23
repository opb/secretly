VERSION := $(shell git describe --tags --always --dirty="-dirty")
LDFLAGS := '-s -w -X "main.Version=$(VERSION)"'

buildall: clean
	mkdir build
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o build/secretly-$(VERSION)-darwin-amd64 -ldflags=$(LDFLAGS) -v ./cmd
	GOOS=linux  GOARCH=amd64 CGO_ENABLED=0 go build -o build/secretly-$(VERSION)-linux-amd64  -ldflags=$(LDFLAGS) -v ./cmd
clean:
	rm -rf build
