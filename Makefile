GOFILES = main/*.go
build: ${GOFILES}
	go build -o kubectl-repl ${GOFILES}
