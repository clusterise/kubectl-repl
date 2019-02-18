FROM golang:1.11-alpine3.9 as build

ARG VERSION

WORKDIR /go/src/github.com/Mikulas/kubectl-repl
ADD . /go/src/github.com/Mikulas/kubectl-repl/
RUN apk update && \
    apk add --no-cache \
        ca-certificates curl git make && \
    make build VERSION=$VERSION && \
    curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl


FROM alpine:3.9

RUN echo "@edge http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
    echo "@testing http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    apk update && \
    apk add 'readline@edge' 'rlwrap@testing'

VOLUME ["/root/.kube"]

COPY --from=build /go/src/github.com/Mikulas/kubectl-repl/kubectl /usr/bin/kubectl
COPY --from=build /go/src/github.com/Mikulas/kubectl-repl/kubectl-repl /usr/bin/kubectl-repl
COPY ./autocomplete.dic /opt/kubectl.dic
WORKDIR /root
ADD entrypoint.sh .
ENTRYPOINT ["sh", "entrypoint.sh"]
