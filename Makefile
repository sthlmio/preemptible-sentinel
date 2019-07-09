.PHONY: default build clean build-image push-image

BINARY = controller

DOCKER_REPO = eu.gcr.io/sthlmio-public-images/preemptible-sentinel

GIT := $(shell git rev-parse --short HEAD)

GOCMD = go
GOFLAGS ?= $(GOFLAGS:)
LDFLAGS =

default: build

build:
	"$(GOCMD)" build --no-cache ${GOFLAGS} ${LDFLAGS} -o "${BINARY}"

# make build-image VER=0.1.0-alpha.0
build-image:
	@docker build --no-cache -t ${DOCKER_REPO}:${VER} .

# make push-image VER=0.1.0-alpha.0
push-image: build-image
	@docker push ${DOCKER_REPO}:${VER}

# make release VER=0.1.0-alpha.0
release: #push-image
	echo ${VER}
#	sed -i -e "s/^\(\s*version\s*:\s*\).*/\1 $VER/" chart/preemptible-sentinel/Chart.yaml
#	sed -i -e "s/^\(\s*appVersion\s*:\s*\).*/\1 $VER/" chart/preemptible-sentinel/Chart.yaml
#	@docker push ${DOCKER_REPO}:${VER}
#	helm package chart/preemptible-sentinel
#	gsutil cp gs://charts.sthlm.io/index.yaml index.yaml
#	helm repo index --merge index.yaml chart/
#	gsutil cp chart/preemptible-sentinel-${VER}.tgz gs://charts.sthlm.io
#	gsutil cp chart/index.yaml gs://charts.sthlm.io

clean:
	"$(GOCMD)" clean -i