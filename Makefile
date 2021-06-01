BINARY = papyrus
BUILD_DIR = bin
BINARY_PATH = ${BUILD_DIR}/${BINARY}

VERSION=1.0.0
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_VERSION=${VERSION}+${COMMIT}
LDFLAGS = -ldflags "-X github.com/ahodieb/papyrus/cmd.version=${BUILD_VERSION}"

all: clean test build

test: 
	go test -v ./...

build: test
	go build ${LDFLAGS} -o ${BINARY_PATH}

clean:
	-rm -rf ${BUILD_DIR}

.PHONY: clean test build all
