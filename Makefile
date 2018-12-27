IMAGE_NAME := yrashetska/intest
CUR_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

rebuild:
	go build -o run
	./run -suites test
	rm -f run
rebuild-debug:
	go build -o run
	./run -suites test -debug
	rm -f run
run-test:
	go test ./...

# Create the Docker image with the latest tag.
.PHONY: latest
build-latest:
	go build -o run
	docker build -f Dockerfile -t $(IMAGE_NAME):latest .
	rm -f run

push-latest:
	docker push $(IMAGE_NAME):latest

#Explore docker network
#docker network ls
#Use --net some_net
#To get an access to the containers which have to be linked to the test container
docker-test:
	docker run --rm -v $(CURDIR)/test:/usr/local/bin/test --net some-network --link auth:auth --link gw:gw $(IMAGE_NAME) -suites test