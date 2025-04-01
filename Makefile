# Makefile for Go projects.
#
# This Makefile makes an effort to provide standard make targets, as described
# by https://www.gnu.org/prep/standards/html_node/Standard-Targets.html.
SHELL := /bin/sh

include ./rules/Makefile.*

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

################################################################################
## Standard make targets
################################################################################

.DEFAULT_GOAL := all
.PHONY: all
all: setup generate fix install

.PHONY: setup
setup:
	touch .env

.PHONY: install
install: go-install

.PHONY: uninstall
uninstall:
	@echo "Uninstalling ${SETTINGS_PROJECT_NAME}"
	$rm -f $(GOPATH)/bin/${SETTINGS_PROJECT_NAME}

.PHONY: clean
clean: go-clean
	@echo "Deleting coverage"
	@rm -f coverage.out

.PHONY: check
check: test

################################################################################
## Go-like targets
################################################################################

.PHONY: build
build: go-build

.PHONY: generate
generate: go-generate

.PHONY: test
test: go-test

.PHONY: cover
cover: go-cover/text

.PHONY: cover/html
cover/html: go-cover/html

.PHONY: cover/text
cover/text: go-cover/text

################################################################################
## Linters and formatters
################################################################################

.PHONY: fix
fix: go-fix

.PHONY: lint
lint: go-lint

################################################################################
## Docker
################################################################################

.PHONY: image
image:
	@echo "> Building dev docker image"
	docker build -t dexlabsio/garlic:dev .
