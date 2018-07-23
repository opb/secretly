VERSION := $(shell git describe --tags --always --dirty="-dev")
LDFLAGS := -ldflags='-s -w -X "main.Version=$(VERSION)"'
GT := $$GITHUB_TOKEN
TAG := $$CIRCLE_TAG

ifeq ($(TAG),"")
$(info CIRCLE_TAG not set!)
pushrelease:
else
pushrelease: release
endif

release: gh-release dist
	github-release release \
	--security-token $(GT) \
	--user opb \
	--repo secretly \
	--tag $(VERSION) \
	--name $(VERSION)

	github-release upload \
	--security-token $(GT) \
	--user opb \
	--repo secretly \
	--tag $(VERSION) \
	--name secretly-$(VERSION)-darwin-amd64 \
	--file dist/secretly-$(VERSION)-darwin-amd64

	github-release upload \
	--security-token $(GT) \
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