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
PROJECT     := vamp-workflow-agent
SRCDIR      := $(CURDIR)
DESTDIR     := target

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
export CGO_ENABLED := 0
export LDFLAGS     := "-X main.version=$(VERSION)"
export GOFLAGS     := -a -installsuffix cgo



# Targets
.PHONY: all
all: default

# Using our buildserver which contains all the necessary dependencies
.PHONY: default
default:
	docker pull $(BUILD_SERVER)
	docker run \
		--rm \
		--volume /var/run/docker.sock:/var/run/docker.sock \
		--volume $(shell command -v docker):/usr/bin/docker \
		--volume $(CURDIR):/srv/src/go/src/github.com/magneticio/vamp-workflow-agent \
		--workdir=/srv/src/go/src/github.com/magneticio/vamp-workflow-agent \
		$(BUILD_SERVER) \
			make $(PROJECT)

$(PROJECT):
	@echo "Building: $(PROJECT)_$(VERSION)_$(GOOS)_$(GOARCH)"
	mkdir -p $(DESTDIR)/vamp
	go get -d ./...
	go install
	go build -ldflags $(LDFLAGS) $(GOFLAGS) -o $(DESTDIR)/vamp/$(PROJECT)


.PHONY: build_npm
build_npm:
	@echo "Installing vamp-node-client"
	mkdir -p $(DESTDIR)/vamp
	npm install --prefix $(DESTDIR)/vamp git://github.com/magneticio/vamp-node-client
	npm install --prefix /tmp removeNPMAbsolutePaths
	/tmp/node_modules/.bin/removeNPMAbsolutePaths $(DESTDIR)/vamp

.PHONY: build_ui
build_ui:
	@echo "Building ui"
	npm install --prefix $(SRCDIR)/ui
	cd $(SRCDIR)/ui ; $(SRCDIR)/ui/node_modules/.bin/ng build --env=prod
	mv $(SRCDIR)/ui/dist $(DESTDIR)/vamp/ui

.PHONY: docker
docker: $(PROJECT) build_npm
	mkdir -p $(DESTDIR)/docker
	cp $(SRCDIR)/Dockerfile $(DESTDIR)/docker/Dockerfile
	cp -Rf $(SRCDIR)/files $(DESTDIR)/docker
	tar -C $(DESTDIR) -zcvf $(PROJECT)_$(VERSION)_$(GOOS)_$(GOARCH).tar.gz vamp
	mv $(PROJECT)_$(VERSION)_$(GOOS)_$(GOARCH).tar.gz $(DESTDIR)/docker

.PHONY: clean
clean: clean-$(PROJECT) clean-docker

.PHONY: clean-$(PROJECT)
clean-$(PROJECT):
	rm -rf $(DESTDIR)/vamp/$(PROJECT)

.PHONY: clean-docker
clean-docker:
	rm -rf $(DESTDIR)/docker
