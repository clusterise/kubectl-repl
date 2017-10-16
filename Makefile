GOFILES = main/*.go
VERSION ?= master
REPO = mikulas/kubectl-repl

build: ${GOFILES} kubectl-repl

kubectl-repl: ${GOFILES}
	cd main && go get -t ./... && go install
	go build -o kubectl-repl ${GOFILES}

preflight: format test

format:
	go fmt ${GOFILES}

test:
	go test ${GOFILES}

release:
	bash release.sh ${VERSION}

docker:
	cd docker && docker build -t ${REPO}:${VERSION} .

.PHONY: docker preflight
