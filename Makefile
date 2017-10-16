GOFILES = main/*.go
VERSION ?= master

build: ${GOFILES}
	cd main && go install
	go build -o kubectl-repl ${GOFILES}

release:
	bash release.sh ${VERSION}
