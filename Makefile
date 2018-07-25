VERSION := $(shell git describe --tags --always --dirty="-dev")
IS_TAG := $(shell git describe --tags --candidates=0 2> /dev/null)
LDFLAGS := -ldflags='-s -w -X "main.Version=$(VERSION)"'
IS_VERSION := $(shell echo $(VERSION) | cut -c1-1)
IS_DEV := $(lastword $(subst -dev, dev,$(VERSION)))

release: gh-release dist
ifeq ($(IS_TAG),)
		@echo "Not tagged with a version, aborting"
		exit 2
endif
ifeq ($(IS_DEV),dev)
		@echo "HEAD is dirty, aborting"
		exit 2
endif
	github-release release \
        --security-token $(GITHUB_TOKEN) \
        --user opb \
        --repo secretly \
        --tag $(VERSION) \
        --name $(VERSION)

	github-release upload \
        --security-token $(GITHUB_TOKEN) \
        --user opb \
        --repo secretly \
        --tag $(VERSION) \
        --name secretly-$(VERSION)-darwin-amd64 \
        --file dist/secretly-$(VERSION)-darwin-amd64

	github-release upload \
        --security-token $(GITHUB_TOKEN) \
        --user opb \
        --repo secretly \
        --tag $(VERSION) \
        --name secretly-$(VERSION)-linux-amd64 \
        --file dist/secretly-$(VERSION)-linux-amd64

dist: clean
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -v -o dist/secretly-$(VERSION)-darwin-amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -v -o dist/secretly-$(VERSION)-linux-amd64

clean:
	rm -rf dist/*

gh-release:
	go get -u github.com/aktau/github-release