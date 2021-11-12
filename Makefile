_CURDIR := `git rev-parse --show-toplevel 2>/dev/null | sed -e 's/(//'`
_STATIC_DIR := "${_CURDIR}/cmd/static/dist"
_UI_DIR := "${_CURDIR}/ui/dist"
_BUILD_DIR := "${_CURDIR}/build"

STATIC_IMAGE := "test-static:latest"
API_IMAGE := "test-api:latest"

STATIC_NAME_CN = "test-static"
API_NAME_CN = "test-api"

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
		--tag "${STATIC_IMAGE}" \
		--file "${_BUILD_DIR}/docker/static.docker" \
		"${_BUILD_DIR}"
	@docker build \
		--tag "${API_IMAGE}" \
		--file "${_BUILD_DIR}/docker/service.docker" \
		"${_BUILD_DIR}"

.PHONY: clean
clean: ## Clean
	@rm -vf "${STATIC_OUT}" "${API_OUT}"
	@rm -rf ${_STATIC_DIR}/*
	@go clean
	@docker rmi "${STATIC_IMAGE}" "${API_IMAGE}"

.PHONY: run
run: ## Run docker services
	@docker stack deploy -c "${_BUILD_DIR}/docker/docker-compose.yml" test

stop: ## Stop docker services
	@docker stack rm test


.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'