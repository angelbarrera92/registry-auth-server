.DEFAULT_GOAL: help
SHELL := /bin/bash

PROJECTNAME := $(shell basename "$(PWD)")

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: build
## build: Builds the registry-auth-server container image.
build:
	@docker build --pull --no-cache -t registry-auth-server:latest --build-arg COMMIT=${GITHUB_SHA} --build-arg VERSION=$(TAG_VERSION) -f build/Dockerfile .

.PHONY: push
## push: Push the registry-auth-server container image. Ensure you are logged in the registry before pushing.
push:
	@docker tag registry-auth-server:latest angelbarrera92/registry-auth-server:$(TAG_VERSION)
	@docker push angelbarrera92/registry-auth-server:$(TAG_VERSION)
