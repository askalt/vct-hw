VERSION ?= 0.0.1
BACKEND_NAME := events-backend
DOCKER_HUB_BACKEND_REPO := askalt/vct-hw-backend

FRONTEND_NAME := events-frontend
DOCKER_HUB_FRONTEND_REPO := askalt/vct-hw-frontend


.PHONY: dist
dist:
	mkdir -p dist


.PHONY: build/backend
build/backend: dist
	go build -C backend -o ../dist/${BACKEND_NAME}


build/docker/backend:
	docker build -t ${BACKEND_NAME}:${VERSION} backend
	docker tag ${BACKEND_NAME}:${VERSION} ${BACKEND_NAME}:latest
	docker tag ${BACKEND_NAME}:${VERSION} ${DOCKER_HUB_BACKEND_REPO}:latest


.PHONY: build/frontend
build/frontend: dist
	go build -C frontend -o ../dist/${FRONTEND_NAME}


build/docker/frontend:
	docker build -t ${FRONTEND_NAME}:${VERSION} frontend
	docker tag ${FRONTEND_NAME}:${VERSION} ${FRONTEND_NAME}:latest
	docker tag ${FRONTEND_NAME}:${VERSION} ${DOCKER_HUB_FRONTEND_REPO}:latest


build: build/backend build/frontend


build/docker: build/docker/backend build/docker/frontend


unit-test:
	go test -v backend
	go test -v frontend


.PHONY: clean
clean:
	rm -rf dist
