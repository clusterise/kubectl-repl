GOFILES = main/*.go
VERSION ?= master
REPO = mikulas/kubectl-repl

build: ${GOFILES} kubectl-repl

kubectl-repl: ${GOFILES}
	cd main && go get -t ./... && go install
	sed -i.bak "s/{{{VERSION}}}/${VERSION}/" main/main.go
	go build -o kubectl-repl ${GOFILES}
	mv main/main.go.bak main/main.go

preflight: format test

format:
	go fmt ${GOFILES}

test:
	go test ${GOFILES}

release:
	sed -i.bak "s/{{{VERSION}}}/${VERSION}/" main/main.go
	bash release.sh ${VERSION}
	mv main/main.go.bak main/main.go

docker:
	rm kubectl-repl || true
	docker build --build-arg VERSION=${VERSION} -t ${REPO}:${VERSION} .

.PHONY: docker preflight
