GOFILES = main/*.go
build: ${GOFILES}
	cd main && go install
	go build -o kubectl-repl ${GOFILES}

release:
	bash release.sh
