#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

VERSION="$1"

function build {
    GOOS="$1" GOARCH="$2" go build -o "releases/kubectl-repl-$1-$2-$VERSION" main/*.go
}

mkdir -p "$DIR/releases"

set -v
build windows 386
build windows amd64
build darwin 386
build darwin amd64
build linux 386
build linux amd64
build linux arm
