#!/usr/bin/env bash
[ "$CIRCLE_TAG" ] || exit 0;

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