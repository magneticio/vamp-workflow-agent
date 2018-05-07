# See: http://clarkgrubb.com/makefile-style-guide
SHELL             := bash
.SHELLFLAGS       := -eu -o pipefail -c
.DEFAULT_GOAL     := default
.DELETE_ON_ERROR  :
.SUFFIXES         :

STASH       := stash
PROJECT     := vamp-workflow-agent
PROJECT_DIR := $(CURDIR)
VERSION     := $(shell git describe --tags)
TARGET      := $(CURDIR)/target
IMAGE_TAG   := $(shell echo $$BRANCH_NAME)

# if Makefile.local exists, include it.
ifneq ("$(wildcard Makefile.local)", "")
	include Makefile.local
endif

ifeq ($(strip $(IMAGE_TAG)),)
IMAGE_TAG := $(shell git rev-parse --abbrev-ref HEAD)
endif

.PHONY: clean
clean:
	rm -rf $(TARGET)

.PHONY: tag
tag:
	@echo $(IMAGE_TAG)

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: frontend
frontend:
	@cd $(CURDIR)/ui && make STASH=$(STASH) PROJECT_DIR=$(PROJECT_DIR) local

.PHONY: backend
backend:
	uname -a
	@go version
	@echo "$(PROJECT): $(VERSION)"
	go get -d -v ./...
	go install -v ./...
	go build -ldflags "-X main.version=$(VERSION)" -a -installsuffix cgo

.PHONY: build
build:
	@cd $(CURDIR)/ui && make purge
	docker run \
         --rm \
         --volume $(STASH):/go \
         --volume $(PROJECT_DIR):/go/src/github.com/magneticio/$(PROJECT) \
         --workdir=/go/src/github.com/magneticio/$(PROJECT) \
         golang:1.10.2 make STASH=$(STASH) PROJECT_DIR=$(PROJECT_DIR)/ui backend
	@cd $(CURDIR)/ui && \
	docker run \
         --rm \
         --volume $(STASH):/root \
         --volume $(PROJECT_DIR)/ui:/$(PROJECT) \
         --workdir=/$(PROJECT) \
         node:9.11.1 make local
	mkdir -p $(TARGET)
	echo $(VERSION) > $(TARGET)/version
	cp $(CURDIR)/Dockerfile $(TARGET)/ && cp -R $(CURDIR)/files $(TARGET)/
	mv $(CURDIR)/$(PROJECT) $(TARGET)/. && mv $(CURDIR)/ui/dist $(TARGET)/ui
	docker build -t magneticio/$(PROJECT):$(IMAGE_TAG) $(TARGET)

.PHONY: default
default: clean build
