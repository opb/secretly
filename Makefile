VERSION := $(shell git describe --tags --always --dirty="-dev")
LDFLAGS := -ldflags='-s -w -X "main.Version=$(VERSION)"'
GT := $$GITHUB_TOKEN
TAG := $$CIRCLE_TAG

release: gh-release dist
	@if [ "$(TAG)" = "" ]; then\
		echo "Not a tag - skipping";\
	else
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
	fi

dist: clean
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -v -o dist/secretly-$(VERSION)-darwin-amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -v -o dist/secretly-$(VERSION)-linux-amd64

clean:
	rm -rf dist/*

gh-release:
	go get -u github.com/aktau/github-release