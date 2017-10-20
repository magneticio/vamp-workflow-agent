# See: http://clarkgrubb.com/makefile-style-guide
SHELL             := bash
.SHELLFLAGS       := -eu -o pipefail -c
.DEFAULT_GOAL     := default
.DELETE_ON_ERROR:
.SUFFIXES:

# Targets
.PHONY: all
all: default

# Install the necessary NodeJS dependencies and build the UI
.PHONY: default
default:
	@echo "Building ui"
	npm install
	npm run ng -- build --env=prod

.PHONY: clean
clean:
	rm -rf $(CURDIR)/node_modules
	rm -rf $(CURDIR)/dist
	rm -rf $(CURDIR)/npm-debug.log
