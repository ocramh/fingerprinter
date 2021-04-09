PKG=github.com/ocramh/fingerprinter

.PHONY: install
install:
	go install ${PKG}/cmd/fingerprinter

.PHONY: docker-run
docker-run:
	docker build -t sygma/fingerprinter .
	docker run -it \
		--rm \
		--entrypoint=/bin/sh \
		--name fingerprinter \
		sygma/fingerprinter