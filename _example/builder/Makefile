NAME=builder
VERSION=v0.0.1

GIT_IMPORT=github.com/things-go/dyn/core/builder
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_FULL_COMMIT=$(shell git rev-parse HEAD)
GIT_TAG=$(shell git describe --abbrev=0 --tags --always --match "v*")
BUILD_DATE=$(shell date "+%F %T %z")

LDFLAGS=-X '$(GIT_IMPORT).BuildDate=$(BUILD_DATE)' \
		-X '$(GIT_IMPORT).Version=$(VERSION)' \
		-X '$(GIT_IMPORT).GitCommit=$(GIT_COMMIT)' \
		-X '$(GIT_IMPORT).GitFullCommit=$(GIT_FULL_COMMIT)' \
		-X '$(GIT_IMPORT).GitTag=$(GIT_TAG)'

build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) go build -trimpath -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}" -o ${NAME}

clean:
	rm -rf ./${NAME}

.PHONY: build clean