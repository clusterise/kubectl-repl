GOFILES = main/*.go
VERSION ?= master

build: ${GOFILES}
	cd main && go install
	go build -o kubectl-repl ${GOFILES}

test:
	go test ${GOFILES}

release:
	bash release.sh ${VERSION}
