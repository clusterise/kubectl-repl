#!/usr/bin/env sh

CONFIG="/root/.kube/config"

if [[ ! -f "$CONFIG" ]]; then
    echo "$CONFIG not found"
    echo "You most likely did not mount your homedir,"
    echo "start this container with volume:"
    echo "    -v \$HOME/.kube:/root/.kube "
    exit 1
fi

# workaround for rlwrap width=0 https://gitlab.com/perl6/docker/issues/1
sleep 0.1

rlwrap kubectl-repl -verbose
