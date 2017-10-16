GOFILES = main/*.go
VERSION ?= master

build: ${GOFILES}
	cd main && go install
	go build -o kubectl-repl ${GOFILES}

preflight: format test

format:
	go fmt ${GOFILES}

test:
	go test ${GOFILES}

release:
	bash release.sh ${VERSION}
