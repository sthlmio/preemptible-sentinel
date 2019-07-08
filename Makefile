.PHONY: default build clean build-image push-image

BINARY = controller

DOCKER_REPO = eu.gcr.io/sthlmio/pvm-controller

GIT := $(shell git rev-parse --short HEAD)

GOCMD = go
GOFLAGS ?= $(GOFLAGS:)
LDFLAGS =

default: build

build:
	"$(GOCMD)" build --no-cache ${GOFLAGS} ${LDFLAGS} -o "${BINARY}"

# make build-image VER=1
build-image:
	@docker build --no-cache -t ${DOCKER_REPO}:${VER}-${GIT} .

# make push-image VER=1
push-image: build-image
	@docker push ${DOCKER_REPO}:${VER}-${GIT}

clean:
	"$(GOCMD)" clean -i