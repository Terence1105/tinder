.phony:

build-image:
	docker build . -f docker/tinder/dockerfile -t tinder:latest

run:
	docker-compose -f docker/local/docker-compose.yml up -d
