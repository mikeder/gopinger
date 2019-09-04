build: 
	@docker build . -t mikeder/gopinger:latest
push:
	@docker push mikeder/gopinger
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
