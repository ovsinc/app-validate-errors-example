_CURDIR := `git rev-parse --show-toplevel 2>/dev/null | sed -e 's/(//'`
_STATIC_DIR := "${_CURDIR}/cmd/static/dist"
_UI_DIR := "${_CURDIR}/ui/dist"
_BUILD_DIR := "${_CURDIR}/build"


REGISTRY_URL = "127.0.0.1:5000"

STATIC_IMG := "test-static"
STATIC_IMG_TAG := "v1"
STATIC_IMAGE := "${STATIC_IMG}:${STATIC_IMG_TAG}"

API_IMG := "test-api"
API_IMG_TAG := "v1"
API_IMAGE := "${API_IMG}:${API_IMG_TAG}"


STATIC_OUT := "${_CURDIR}/build/static"
API_OUT := "${_CURDIR}/build/service"




.PHONY: all
all: build ## Build all (default)

.PHONY: build
build: build_ui build_static build_api build_docker ## Build all

.PHONY: build_ui
build_ui: ## Build the UI (SPA)
	@pushd ${_CURDIR}/ui && quasar build --mode spa; popd


.PHONY: build_static
build_static: ## Build the static server
	@rm -rf ${_STATIC_DIR}/* && \
		cp -r "${_UI_DIR}"/* "${_STATIC_DIR}"
	@CGO_ENABLED=0 \
		go build \
		-v -mod=vendor -installsuffix cgo \
		-o ${STATIC_OUT} \
		${_CURDIR}/cmd/static

.PHONY: build_api
build_api: ## Build the API server
	@CGO_ENABLED=0 \
		go build \
		-v -mod=vendor -installsuffix cgo \
		-o ${API_OUT} \
		${_CURDIR}/cmd/service

.PHONY: build_docker
build_docker: ## Build the Docker images
	@docker build \
		--force-rm \
		--no-cache \
		--tag "${STATIC_IMAGE}" \
		--file "${_CURDIR}/docker/static.docker" \
		"${_BUILD_DIR}"
	@docker tag "${STATIC_IMAGE}" "${STATIC_IMG}:latest"
	@docker tag "${STATIC_IMAGE}" "${REGISTRY_URL}/${STATIC_IMAGE}"
	@docker tag "${STATIC_IMAGE}" "${REGISTRY_URL}/${STATIC_IMG}:latest"
	docker push "${REGISTRY_URL}/${STATIC_IMAGE}"
	docker push "${REGISTRY_URL}/${STATIC_IMG}:latest"
	@docker build \
		--force-rm \
		--tag "${API_IMAGE}" \
		--file "${_CURDIR}/docker/service.docker" \
		"${_BUILD_DIR}"
	@docker tag "${API_IMAGE}" "${API_IMG}:latest"
	@docker tag "${API_IMAGE}" "${REGISTRY_URL}/${API_IMG}:latest"
	@docker tag "${API_IMAGE}" "${REGISTRY_URL}/${API_IMAGE}"
	docker push "${REGISTRY_URL}/${API_IMAGE}"
	docker push "${REGISTRY_URL}/${API_IMG}:latest"

.PHONY: clean
clean: ## Clean
	@rm -vf "${STATIC_OUT}" "${API_OUT}"
	@rm -rf ${_STATIC_DIR}/*
	@go clean
	@docker rmi "${STATIC_IMAGE}" "${API_IMAGE}"

.PHONY: run
run: ## Run docker services
	@docker stack deploy -c "${_CURDIR}/docker/docker-compose.yml" test

.PHONY: start
start: run ## Run docker services

stop: ## Stop docker services
	@docker stack rm test


.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'