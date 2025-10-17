VERSION ?= 0.0.1
APP_NAME := event-server


.PHONY: dist
dist:
	mkdir -p dist


.PHONY: build
build/binary: dist
	go build -o dist/${APP_NAME} .


build/docker-image: build/binary
	docker build -t ${APP_NAME}:${VERSION} .
	docker tag ${APP_NAME}:${VERSION} ${APP_NAME}:latest


unit-test:
	go test -v .


.PHONY: clean
clean:
	rm -rf dist
