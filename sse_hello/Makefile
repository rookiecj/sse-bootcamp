# .env 파일이 있으면 로드
ifneq (,$(wildcard .env))
    include .env
    export
endif

DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' Makefile  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: build
build: 	## Build the project
	go build -o bin/app main.go

.PHONY: run
run: ## Run the project
	go run main.go

.PHONY: clean
clean: ## Clean the project
	rm -f bin/app

.PHONY: test
test: ## Run the tests
	go test -v ./...

.PHONY: lint
lint: ## Run the lint
	go fmt ./...
	go vet ./...

.PHONY: format
format: ## Format the code
	go fmt ./...

# 환경 변수가 설정되어 있지 않은 경우 기본값 사용
REGISTRY_HOST ?= your_registry_host
REGISTRY_USERNAME ?= your_ci_bot_name
REGISTRY_PASSWORD ?= your_ci_bot_password

IMAGE_TAG ?= 1.6
URL_PATH_PREFIX ?= app/sse-hello

.PHONY: docker-setup-builder
docker-setup-builder: ## Setup the docker builder
	docker buildx create --use --name multi-platform-builder || true

.PHONY: docker-build
docker-build: ## Build the docker image for current platform(arm64)
	docker buildx create --use --name multi-platform-builder || true
	# --load 옵션은 단일 플랫폼 빌드에서만 사용 가능합니다.
	docker buildx build --platform linux/arm64 -t sse_hello:$(IMAGE_TAG) --load .

.PHONY: docker-build-multi
docker-build-multi: ## Build the docker image for multiple platforms(amd64,arm64)
	docker buildx create --use --name multi-platform-builder || true
	# 멀티 플랫폼 빌드의 경우 --push 옵션을 사용하여 레지스트리에 직접 푸시해야 합니다.
	docker buildx build --platform linux/amd64,linux/arm64 -t sse_hello:$(IMAGE_TAG) --push .

.PHONY: docker-run
docker-run: ## Run the docker container
	# docker run -p 8080:8080 sse_hello:$(IMAGE_TAG)
	docker run -p 8080:8080 -e URL_PATH_PREFIX=$(URL_PATH_PREFIX) sse_hello:$(IMAGE_TAG)

.PHONY: docker-stop
docker-stop: ## Stop the docker container
	-docker stop sse_hello

.PHONY: docker-remove
docker-remove: docker-stop ## Remove the docker container and image
	-docker rm $(shell docker ps -a -q --filter ancestor=sse_hello:$(IMAGE_TAG))
	-docker rmi -f sse_hello:$(IMAGE_TAG)
	-docker rmi -f $(REGISTRY_HOST)/homelab/sse_hello:$(IMAGE_TAG)

.PHONY: docker-images
docker-images: ## List the docker images
	docker images | grep sse_hello

.PHONY: docker-login
docker-login: ## Login to the docker registry
	@test -n "$(REGISTRY_HOST)" || (echo "REGISTRY_HOST is not set" && exit 1)
	@test -n "$(REGISTRY_USERNAME)" || (echo "REGISTRY_USERNAME is not set" && exit 1)
	@test -n "$(REGISTRY_PASSWORD)" || (echo "REGISTRY_PASSWORD is not set" && exit 1)
	@echo "REGISTRY_HOST: '$(REGISTRY_HOST)'"
	@echo "REGISTRY_USERNAME: '$(REGISTRY_USERNAME)'"
	# @echo "REGISTRY_PASSWORD: '$(REGISTRY_PASSWORD)'"
	@echo '$(REGISTRY_PASSWORD)' | docker login '$(REGISTRY_HOST)' --username '$(REGISTRY_USERNAME)' --password-stdin

.PHONY: docker-tag
docker-tag: ## Tag the docker image
	@test -n "$(REGISTRY_HOST)" || (echo "REGISTRY_HOST is not set" && exit 1)
	docker tag sse_hello:$(IMAGE_TAG) $(REGISTRY_HOST)/homelab/sse_hello:$(IMAGE_TAG)

.PHONY: docker-push
docker-push: docker-tag ## Push the docker image to the registry
	@test -n "$(REGISTRY_HOST)" || (echo "REGISTRY_HOST is not set" && exit 1)
	docker push $(REGISTRY_HOST)/homelab/sse_hello:$(IMAGE_TAG)

.PHONY: docker-pull
docker-pull: ## Pull the docker image from the registry
	@test -n "$(REGISTRY_HOST)" || (echo "REGISTRY_HOST is not set" && exit 1)
	docker pull $(REGISTRY_HOST)/homelab/sse_hello:$(IMAGE_TAG)

.PHONY: docker-pull-latest
docker-pull-latest: ## Pull the docker image from the registry
	[ -z "$(REGISTRY_HOST)" ] && echo "REGISTRY_HOST is not set" && exit 1
	docker pull $(REGISTRY_HOST)/homelab/sse_hello:latest

.PHONY: generate-certs
generate-certs: ## Generate SSL certificates for Harbor
	mkdir -p certs
	cd certs && ./generate-certs.sh

.PHONY: generate-favicon
generate-favicon: ## Generate favicon from SVG
	which magick || brew install imagemagick
	chmod +x scripts/generate-favicon.sh
	./scripts/generate-favicon.sh

