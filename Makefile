.PHONY: default build docker-image test stop clean-images clean

BINARY = kubestatewatch

VERSION= $(shell git describe --tags --always --dirty)
BUILD= $(shell date +%FT%T%z)

PKG            = github.com/marvasgit/kubestatewatch
TRAVIS_COMMIT ?= `git describe --tags`
GOCMD          = go
BUILD_DATE     = `date +%FT%T%z`
GOFLAGS       ?= $(GOFLAGS:)
LDFLAGS       := "-X '$(PKG)/cmd.gitCommit=$(TRAVIS_COMMIT)' \
		          -X '$(PKG)/cmd.buildDate=$(BUILD_DATE)'"
BUILD_VERSION ?= $(shell cat VERSION)
DOCKER_IMAGE ?= docker.io/teadove/${BINARY}:$(BUILD_VERSION)

default: build test

build:
	"$(GOCMD)" build ${GOFLAGS} -ldflags ${LDFLAGS} -o "${BINARY}"

docker-image:
	@docker buildx build --platform linux/amd64 . --tag "${DOCKER_IMAGE}" --push

test:
	"$(GOCMD)" test -race -v ./...

stop:
	@docker stop "${BINARY}"

clean-images: stop
	@docker rmi "${BUILDER}" "${BINARY}"

clean:
	"$(GOCMD)" clean -i
