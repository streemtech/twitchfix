
.PHONY: docker main

main:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o twitchfix .

docker: main
	docker build \
	--no-cache \
	-t twitchfix \
	.

deploy: docker
	docker run \
	--name twitchfix \
	-p 8284:8284 \
	twitchfix
