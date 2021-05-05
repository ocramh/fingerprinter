PKG=github.com/ocramh/fingerprinter
DOCKER_IMAGE_NAME=ocramh/fingerprinter
COVERAGE_DIR=coverage

.DEFAULT_GOAL = help

.PHONY: help
help:  ## shows this help message
	@echo 'usage: make [target] ...'
	@echo
	@echo 'targets:'
	@echo
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s : | sed -e 's/^/  /')"

.PHONY: install
install: ## installs the executable
	go install ${PKG}/cmd/fingerprinter

.PHONY: test
test: ## runs unit tests
	@mkdir -p $(COVERAGE_DIR)
	@echo 'mode: atomic' > $(COVERAGE_DIR)/coverage.out
	@go test ./... -coverprofile=$(COVERAGE_DIR)/coverage.out
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

.PHONY: docker-run
docker-run: ## builds and runs the app docker container
	docker build -t $(DOCKER_IMAGE_NAME) .
	docker run -it \
		--rm \
		--entrypoint=/bin/sh \
		--name fingerprinter \
		$(DOCKER_IMAGE_NAME)