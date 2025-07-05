DOCKERCMD=docker

DOCKER_CONTAINER_NAME?=koer-tax-service
DOCKER_CONTAINER_IMAGE?=koer-tax-service:latest
DOCKER_BUILD_ARGS?=

BUILD_DATE?=$(shell date -u +'%Y-%m-%dT00:00:00Z')
BUILD_VERSION?=1.0.0

TOPDIR=$(PWD)
BINARY=koer-tax-service
.FORCE:
.PHONY: build
.PHONY: vet
.PHONY: unit-test
.PHONY: generates
.PHONY: depend
.PHONY: docker-build
.PHONY: solr
.PHONY: clean
.PHONY: install
.PHONY: install-wkhtmltopdf
.PHONY: all
.PHONY: .FORCE

guard-%:
	@if [ -z '${${*}}' ]; then echo 'Environment variable $* not set' && exit 1; fi

build:
	@echo "Executing go build"
	go build -v -buildmode=pie -ldflags "-X main.version=$(BUILD_VERSION)" -o app ./cmd/api/
	@echo "Binary ready"

vet:
	@echo "Running Go static code analysis with go vet"
	go vet -asmdecl -atomic -bool -buildtags -copylocks -httpresponse -loopclosure -lostcancel -methods -nilfunc -printf -rangeloops -shift -structtag -tests -unreachable -unsafeptr ./...
	@echo "go vet complete"

unit-test:
	@echo "Executing go unit test"
	go test -v -cover -json -count=1 -parallel=4 ./...
	@echo "Unit test done"

generate:
	go generate ./...

run:
	go run ./cmd/api/ grpc-gw-server --port1 9224 --port2 3224 --grpc-endpoint :9224 

migrate-db:
	go run ./cmd/api/ db-migrate

depend:
	@echo "Pulling all Go dependencies"
	go mod download
	go mod verify
	go mod tidy
	go mod vendor
	@echo "You can now run 'make build' to compile all packages"

docker-build:
	$(DOCKERCMD) build -t $(DOCKER_CONTAINER_IMAGE) --build-arg GOPROXY=$(GOPROXY) --build-arg GOSUMDB=$(GOSUMDB) --build-arg BUILD_VERSION=$(BUILD_VERSION) $(DOCKER_BUILD_ARGS) .

default: depend

all: depend install-wkhtmltopdf generate build unit-test

install: depend install-wkhtmltopdf build

clean:
	rm -f $(BINARY)
	rm -f $(BINARY).exe

