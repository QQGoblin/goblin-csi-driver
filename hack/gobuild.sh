#!/bin/bash

export HTTPS_PROXY=http://172.24.37.57:8001
export HTTP_PROXY=http://172.24.37.57:8001

GOOS=linux GOARCH=amd64 CGO_ENABLED=0  GO111MODULE=on GOPROXY=https://goproxy.cn go mod tidy
GOOS=linux GOARCH=amd64 CGO_ENABLED=0  GO111MODULE=on GOPROXY=https://goproxy.cn go mod vendor
GOOS=linux GOARCH=amd64 CGO_ENABLED=0  GO111MODULE=on GOPROXY=https://goproxy.cn go build -o ./output/hostpath-csi-driver ./cmd/hostpath/main.go