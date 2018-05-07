# See: http://clarkgrubb.com/makefile-style-guide
SHELL             := bash
.SHELLFLAGS       := -eu -o pipefail -c
.DEFAULT_GOAL     := default
.DELETE_ON_ERROR  :
.SUFFIXES         :

STASH       := stash
PROJECT     := vamp-workflow-agent
PROJECT_DIR := $(CURDIR)

# if Makefile.local exists, include it.
ifneq ("$(wildcard Makefile.local)", "")
	include Makefile.local
endif

.PHONY: clean
clean:
	rm -rf "$(CURDIR)"/dist

.PHONY: purge
purge: clean
	rm -rf "$(CURDIR)"/node_modules

.PHONY: update
update:
	yarn install
	yarn upgrade

.PHONY: local
local:
	uname -a
	yarn install
	yarn run ng build -prod

.PHONY: build
build:
	docker run \
         --rm \
         --volume $(STASH):/root \
         --volume $(PROJECT_DIR):/$(PROJECT) \
         --workdir=/$(PROJECT) \
         node:9.11.1 make local

.PHONY: default
default: clean build
