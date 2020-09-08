
.PHONY: docker main dockerhub deploy standalone run

main:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o twitchfix .

#to run
#make VERSION=0.0.2 dockerhub

docker: main
	docker build \
	--no-cache \
	-t streemtech/twitchfix:$(VERSION) \
	.

dockerhub: docker
	docker push streemtech/twitchfix:$(VERSION)

deploy: docker
	docker run \
	--name streemtech/twitchfix:$(VERSION) \
	-p 8284:8284 \
	twitchfix

standalone: 
	go build -o twitchfix .

run: standalone
	./twitchfix