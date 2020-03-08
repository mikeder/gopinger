DEFAULT_TARGET=build
GOFLAGS=-mod=vendor

build:
	@go build -mod vendor -o ./bin/gopinger ./cmd/
.PHONY: build

build-image:
	@docker build . -t mikeder/gopinger:latest
.PHONY: build-image

push:
	@docker push mikeder/gopinger
.PHONY: push

revendor:
	@go mod tidy -v
	@go mod vendor -v
	@go mod verify
	@git add -A vendor
.PHONY: revendor

update:
	@go get -u -mod=''
	@make revendor
.PHONY: update
