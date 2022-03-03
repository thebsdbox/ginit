SHELL := /bin/sh

# The name of the executable (default is current directory name)
TARGET := ginit
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 0.0.1
BUILD := `git rev-parse HEAD`

# Operating System Default (LINUX)
TARGETOS=linux

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-s -w -X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -extldflags -static"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

DOCKERTAG ?= $(VERSION)
REPOSITORY = thebsdbox

.PHONY: all build clean install uninstall fmt simplify check run

all: check install

$(TARGET): $(SRC)
	@go build $(LDFLAGS) -o $(TARGET)

build: $(TARGET)
	@true

clean:
	@rm -f $(TARGET)

install:
	@echo Building and Installing project
	@go install $(LDFLAGS)

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

demo:
	@cd demo
	@docker buildx build  --platform linux/amd64,linux/arm64,linux/arm/v7 --push -t $(REPOSITORY)/$(TARGET):$(DOCKERTAG) .
	@echo New Multi Architecture Docker image created
	@cd ..

# This is typically only for quick testing
dockerx86:
	@docker buildx build  --platform linux/amd64 --load -t $(REPOSITORY)/$(TARGET):$(DOCKERTAG) -f initrd.Dockerfile .
	@echo New Multi Architecture Docker image created

docker:
	@docker buildx build  --platform linux/amd64,linux/arm64,linux/arm/v7 --push -t $(REPOSITORY)/$(TARGET):$(DOCKERTAG) -f initrd.Dockerfile .
	@echo New Multi Architecture Docker image created

# This is typically only for quick testing
getramdisk:

	@ID=$$(docker create $(REPOSITORY)/$(TARGET):$(DOCKERTAG) null); \
	docker cp $$ID:/initramfs.cpio.gz initramfs.cpio.gz ; docker rm $$ID
	@echo Extracted ramdisk

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done
	@go tool vet ${SRC}

run: install
	@$(TARGET)