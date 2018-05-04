# See: http://clarkgrubb.com/makefile-style-guide
SHELL             := bash
.SHELLFLAGS       := -eu -o pipefail -c
.DEFAULT_GOAL     := default
.DELETE_ON_ERROR  :
.SUFFIXES         :

STASH       := stash
PROJECT     := vamp-workflow-agent
PROJECT_DIR := $(CURDIR)
FABRICATOR  := golang:1.10.2
VERSION     := $(shell git describe --tags)
TARGET      := $(CURDIR)/target

# if Makefile.local exists, include it.
ifneq ("$(wildcard Makefile.local)", "")
	include Makefile.local
endif

.PHONY: clean
clean:
	rm -rf $(TARGET)

.PHONY: local
local:
	@go version
	@echo "$(PROJECT): $(VERSION)"
	go get -d -v ./...
	go install -v ./...
	go build -ldflags "-X main.version=$(VERSION)" -a -installsuffix cgo
	@cd $(CURDIR)/ui && make STASH=$(STASH) PROJECT_DIR=$(PROJECT_DIR) build
	mkdir -p $(TARGET)
	mv $(CURDIR)/$(PROJECT) $(TARGET)/.
	mv $(CURDIR)/ui/dist $(TARGET)/ui
	cp $(CURDIR)/Dockerfile $(TARGET)/
	echo $(VERSION) > $(TARGET)/version
	cp -R $(CURDIR)/files $(TARGET)/
	docker build -t magneticio/$(PROJECT):$$(git rev-parse --abbrev-ref HEAD) $(TARGET)

.PHONY: build
build:
	docker run \
         --rm \
         --volume $(STASH):/go \
         --volume $(PROJECT_DIR):/go/src/github.com/magneticio/$(PROJECT) \
         --volume $$(which docker):/usr/bin/docker \
         --volume /var/run/docker.sock:/var/run/docker.sock \
         --workdir=/go/src/github.com/magneticio/$(PROJECT) \
         $(FABRICATOR) make STASH=$(STASH) PROJECT_DIR=$(PROJECT_DIR)/ui local

.PHONY: default
default: clean build
