build: deps
	@docker build . -t mikeder/gopinger:latest
deps:
	@cp /etc/ssl/certs/ca-certificates.crt .
push:
	@docker push mikeder/gopinger
