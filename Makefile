# See: http://clarkgrubb.com/makefile-style-guide
SHELL             := bash
.SHELLFLAGS       := -eu -o pipefail -c
.DEFAULT_GOAL     := default
.DELETE_ON_ERROR:
.SUFFIXES:

# Constants, these can be overwritten in your Makefile.local
BUILD_SERVER := magneticio/buildserver

# if Makefile.local exists, include it.
ifneq ("$(wildcard Makefile.local)", "")
	include Makefile.local
endif

# Directories
PROJECT := vamp-workflow-agent
SRCDIR  := $(CURDIR)
DESTDIR := target

# Determine which version we're building
ifeq ($(shell git describe --tags),$(shell git describe --abbrev=0 --tags))
	export VERSION := $(shell git describe --tags)
else
	export VERSION := katana
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
export CGO_ENABLED  := 0
export LDFLAGS      := "-X main.version=$(VERSION)"
export GOFLAGS      := -a

# Using our buildserver which contains all the necessary dependencies
.PHONY: default
default:
	docker pull $(BUILD_SERVER)
	docker run \
		--rm \
		--volume $(CURDIR):/srv/src/go/src/github.com/magneticio/vamp-workflow-agent \
		--workdir=/srv/src/go/src/github.com/magneticio/vamp-workflow-agent \
		$(BUILD_SERVER) \
			make docker-context

	$(MAKE) docker

# Build the 'vamp-workflow-agent' go binary
$(PROJECT):
	@echo "Building: $(PROJECT)_$(VERSION)_$(GOOS)_$(GOARCH)"
	mkdir -p $(DESTDIR)/vamp
	go get -d ./...
	go build -ldflags $(LDFLAGS) $(GOFLAGS) -o $(DESTDIR)/vamp/$(PROJECT)

# Install the necessary NodeJS dependencies
.PHONY: build-npm
build-npm:
	@echo "Installing vamp-node-client"
	mkdir -p $(DESTDIR)/vamp
	npm install --prefix $(DESTDIR)/vamp git://github.com/magneticio/vamp-node-client
	npm install --prefix /tmp removeNPMAbsolutePaths
	/tmp/node_modules/.bin/removeNPMAbsolutePaths $(DESTDIR)/vamp

# Build the UI
# All UI build steps are managed in a separate Makefile in the 'ui' directory
.PHONY: build-ui
build-ui:
	@echo "Building ui"
	$(MAKE) -C $(SRCDIR)/ui
	[[ -d $(DESTDIR)/vamp/ui ]] && (rm -rf $(DESTDIR)/vamp/ui)
	mv $(SRCDIR)/ui/dist $(DESTDIR)/vamp/ui

# Copying all necessary files and setting version under 'target/docker/'
.PHONY: docker-context
docker-context: $(PROJECT) build-npm build-ui
	@echo "Creating docker build context"
	mkdir -p $(DESTDIR)/docker
	cp $(SRCDIR)/Dockerfile $(DESTDIR)/docker/Dockerfile
	cp -Rf $(SRCDIR)/scripts $(DESTDIR)/docker
	tar -C $(DESTDIR) -zcvf $(PROJECT)_$(VERSION)_$(GOOS)_$(GOARCH).tar.gz vamp
	mv $(PROJECT)_$(VERSION)_$(GOOS)_$(GOARCH).tar.gz $(DESTDIR)/docker
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
clean: clean-$(PROJECT) clean-docker-context clean-ui
	rm -rf $(DESTDIR)/vamp
	rm -rf $(DESTDIR)/docker

.PHONY: clean-$(PROJECT)
clean-$(PROJECT):
	rm -rf $(DESTDIR)/vamp/$(PROJECT)

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
