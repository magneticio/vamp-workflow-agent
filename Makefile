# See: http://clarkgrubb.com/makefile-style-guide
SHELL             := bash
.SHELLFLAGS       := -eu -o pipefail -c
.DEFAULT_GOAL     := default
.DELETE_ON_ERROR:
.SUFFIXES:

# Constants, these can be overwritten in your Makefile.local
BUILD_SERVER := magneticio/buildserver
DIR_NPM	     := $(HOME)/.npm
DIR_GYP	     := $(HOME)/.node-gyp

# if Makefile.local exists, include it.
ifneq ("$(wildcard Makefile.local)", "")
	include Makefile.local
endif

# Directories
BINARY  := vamp-workflow-agent
PROJECT := $(BINARY)
SRCDIR  := $(CURDIR)
DESTDIR := target

# Determine which version we're building
ifeq ($(shell git describe --tags),$(shell git describe --abbrev=0 --tags))
	export VERSION := $(shell git describe --tags)
else
	ifeq ($(VAMP_GIT_BRANCH), $(filter $(VAMP_GIT_BRANCH), master ""))
		export VERSION := katana
	else
		export VERSION := $(subst /,_,$(VAMP_GIT_BRANCH))
	endif
endif

# Determine operating system
OS := $(shell uname -s)
ifeq ($(OS),Darwin)
	export GOOS := darwin
else ifeq ($(OS),Linux)
	export GOOS := linux
else
	export GOOS := linux
endif

# Determine architecture
MACHINE := $(shell uname -m)
ifeq ($(MACHINE),x86_64)
	export GOARCH := amd64
else
	export GOARCH := 386
endif

# Compiler flags
export CGO_ENABLED := 0
export LDFLAGS     := "-X main.version=$(VERSION)"
export GOFLAGS     := -a -installsuffix cgo

# Targets
.PHONY: all
all: default

# Using our buildserver which contains all the necessary dependencies
.PHONY: default
default: clean-check
	docker pull $(BUILD_SERVER)
	docker run \
		--rm \
		--volume $(CURDIR):/srv/src \
		--volume $(DIR_NPM):/home/vamp/.npm \
		--volume $(DIR_GYP):/home/vamp/.node-gyp \
		--workdir=/srv/src \
		--env BUILD_UID=$(shell id -u) \
		--env BUILD_GID=$(shell id -g) \
		$(BUILD_SERVER) \
			make docker-context

	$(MAKE) docker

# Build the 'vamp-workflow-agent' go binary
$(BINARY):
	@echo "Building: $(BINARY)_$(VERSION)_$(GOOS)_$(GOARCH)"
	mkdir -p $(DESTDIR)/vamp
	rm -rf $(DESTDIR)/go/src/github.com/magneticio/$(BINARY)
	mkdir -p $(DESTDIR)/go/src/github.com/magneticio/$(BINARY)
	cp -a *.go $(DESTDIR)/go/src/github.com/magneticio/$(BINARY)
	export GOPATH=$(abspath $(DESTDIR))/go && \
		cd $(DESTDIR)/go/src/github.com/magneticio/$(BINARY) && \
		go get -d ./... && \
		go build -ldflags $(LDFLAGS) $(GOFLAGS) -o $(DESTDIR)/vamp/$(BINARY)

# Install the necessary NodeJS dependencies
.PHONY: build-npm
build-npm:
	@echo "Installing vamp-node-client"
	mkdir -p $(DESTDIR)/vamp
	npm install --prefix $(DESTDIR)/vamp git://github.com/magneticio/vamp-node-client
	npm install --prefix /tmp removeNPMAbsolutePaths@0.0.3
	/tmp/node_modules/.bin/removeNPMAbsolutePaths $(DESTDIR)/vamp

# Build the UI
# All UI build steps are managed in a separate Makefile in the 'ui' directory
.PHONY: build-ui
build-ui:
	@echo "Building ui"
	$(MAKE) -C $(SRCDIR)/ui
	[ -d $(SRCDIR)/ui/dist ] && rm -rf $(DESTDIR)/vamp/ui
	mv $(SRCDIR)/ui/dist $(DESTDIR)/vamp/ui

# Copying all necessary files and setting version under 'target/docker/'
.PHONY: docker-context
docker-context: $(BINARY) build-npm build-ui
	@echo "Creating docker build context"
	mkdir -p $(DESTDIR)/docker
	cp $(SRCDIR)/Dockerfile $(DESTDIR)/docker/Dockerfile
	cp -Rf $(SRCDIR)/files $(DESTDIR)/docker
	tar -C $(DESTDIR) -zcvf $(BINARY)_$(VERSION)_$(GOOS)_$(GOARCH).tar.gz vamp
	mv $(BINARY)_$(VERSION)_$(GOOS)_$(GOARCH).tar.gz $(DESTDIR)/docker
	echo $(VERSION) $$(git describe --tags) > $(DESTDIR)/docker/version

# Building the docker container using the generated context from the
# 'docker-context' target
.PHONY: docker
docker:
	docker build \
		--tag=magneticio/$(PROJECT):$(VERSION) \
		--file=$(DESTDIR)/docker/Dockerfile \
		$(DESTDIR)/docker

# Remove all files copied/generated from the other targets
.PHONY: clean
clean: clean-$(BINARY) clean-docker-context clean-ui
	rm -rf $(DESTDIR)

.PHONY: clean-$(BINARY)
clean-$(BINARY):
	rm -rf $(DESTDIR)/vamp/$(BINARY) $(DESTDIR)/go

.PHONY: clean-docker-context
clean-docker-context:
	rm -rf $(DESTDIR)/docker

.PHONY: clean-ui
clean-ui:
	$(MAKE) -C $(SRCDIR)/ui clean

# Remove the docker image from the system
.PHONY: clean-docker
clean-docker:
	docker rmi magneticio/$(PROJECT):$(VERSION)

.PHONY: clean-check
clean-check:
	if [ $$(find -uid 0 -print -quit | wc -l) -eq 1 ]; then \
		docker run \
		--rm \
		--volume $(CURDIR):/srv/src \
		--workdir=/srv/src \
		$(BUILD_SERVER) \
			make clean; \
	fi
