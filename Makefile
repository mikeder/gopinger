build: 
	@docker build . -t mikeder/gopinger:latest
push:
	@docker push mikeder/gopinger
