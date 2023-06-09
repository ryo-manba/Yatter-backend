BINARY := yatter-backend-go
MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

PATH := $(PATH):${MAKEFILE_DIR}bin
SHELL := env PATH="$(PATH)" /bin/bash
# for go
export CGO_ENABLED = 0
GOARCH = amd64

# 現在のコミットハッシュを取得する
COMMIT=$(shell git rev-parse HEAD)
# 現在のブランチ名を取得する
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GIT_URL=local-git://

LDFLAGS := -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

build: build-linux

# プラットフォームに依存しないバイナリをビルドする
build-default:
	go build ${LDFLAGS} -o build/${BINARY}

build-linux:
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o build/${BINARY}-linux-${GOARCH} .

# 依存関係の準備のためにmodタスクを呼び出す
prepare: mod

mod:
	go mod download

test:
	go test $(shell go list ${MAKEFILE_DIR}/...)

lint:
	if ! [ -x $(GOPATH)/bin/golangci-lint ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.38.0 ; \
	fi
	golangci-lint run --concurrency 2

vet:
	go vet ./...

clean:
	git clean -f -X app bin build

.PHONY:	test clean
