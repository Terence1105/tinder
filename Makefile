.phony:

build-image:
	docker build . -f docker/tinder/dockerfile -t tinder:latest

run: build-image
	docker-compose -f docker/local/docker-compose.yml up -d

run-swag:
	open http://localhost:8080/swagger/index.html#/
	go run ./cmd/tinder/main.go

build-mockgen:
	docker build -t automockgen -f $(PWD)/docker/utils/mockgen.Dockerfile .

mockgen: build-mockgen
	docker run --rm -v $(PWD):/src -w /src automockgen
