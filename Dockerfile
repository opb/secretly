# build stage
FROM golang:1.11.1-alpine3.8 AS build-env
WORKDIR /app
RUN apk add -U curl git make upx
ADD go.* /app/
RUN go mod download
ADD . .
RUN make builddocker
WORKDIR /app/build
RUN upx -3 secretly

# final stage
FROM alpine:3.8
WORKDIR /app

RUN apk add --update bash ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build-env /app/build/secretly /usr/local/bin/