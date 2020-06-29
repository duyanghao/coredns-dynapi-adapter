## server version
SERVER_VERSION = v0.1.0
## Folder content generated files
BUILD_FOLDER = ./build
PROJECT_URL  = github.com/duyanghao/coredns-dynapi-adapter
## command
GO           = go
GO_VENDOR    = go mod
MKDIR_P      = mkdir -p

## Random Alphanumeric String
SECRET_KEY   = $(shell cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)

## UNAME
UNAME := $(shell uname)

################################################

.PHONY: all
all: build.docs build

.PHONY: pre-build
pre-build:
	$(GO_VENDOR) vendor

.PHONY: build
build: pre-build
	$(MAKE) src.build

.PHONY: clean
clean:
	$(RM) -rf $(BUILD_FOLDER)

## src/ ########################################

.PHONY: src.build
src.build:
	$(MKDIR_P) $(BUILD_FOLDER)/pkg/cmd/coredns-dynapi-adapter/
	GO111MODULE=on $(GO) build -mod=vendor -v -o $(BUILD_FOLDER)/pkg/cmd/coredns-dynapi-adapter/coredns-dynapi-adapter \
	./cmd/...

## dockerfiles/ ########################################

.PHONY: dockerfiles.build
dockerfiles.build:
	docker build --no-cache --rm --tag duyanghao/coredns-dynapi-adapter:$(SERVER_VERSION) -f ./docker/Dockerfile .

## git tag version ########################################

.PHONY: push.tag
push.tag:
	@echo "Current git tag version:"$(SERVER_VERSION)
	git tag $(SERVER_VERSION)
	git push --tags

.PHONY: build.docs
build.docs:
	@echo "Build swagger docs"
	swag init -g pkg/route/routes.go
