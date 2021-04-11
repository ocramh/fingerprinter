PKG=github.com/ocramh/fingerprinter

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

.PHONY: docker-run
docker-run: ## builds and runs the app docker container
	docker build -t sygma/fingerprinter .
	docker run -it \
		--rm \
		--entrypoint=/bin/sh \
		--name fingerprinter \
		sygma/fingerprinter