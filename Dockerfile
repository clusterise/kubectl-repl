FROM golang:1.9-alpine3.6 as build

WORKDIR /go/src/github.com/Mikulas/kubectl-repl
ADD . /go/src/github.com/Mikulas/kubectl-repl/
RUN apk update && \
    apk add --no-cache \
        ca-certificates git make && \
    make build


FROM alpine:3.6

RUN echo "@edge http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
    echo "@testing http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    apk update && \
    apk add 'readline@edge' && \
    apk add 'rlwrap@testing'

VOLUME ["/root/.kube"]

COPY --from=build /go/src/github.com/Mikulas/kubectl-repl/kubectl-repl /usr/bin/kubectl-repl
WORKDIR /root
ADD entrypoint.sh .
ENTRYPOINT ["sh", "entrypoint.sh"]
