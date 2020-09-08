
.PHONY: docker main dockerhub deploy standalone run

main:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o twitchfix .

docker: main
	docker build \
	--no-cache \
	-t streemtech/twitchfix \
	.

dockerhub: docker
	docker push streemtech/twitchfix

deploy: docker
	docker run \
	--name streemtech/twitchfix \
	-p 8284:8284 \
	twitchfix

standalone: 
	go build -o twitchfix .

run: standalone
	./twitchfix