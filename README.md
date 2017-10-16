Kubectl REPL
============

[![Go Report Card](https://goreportcard.com/badge/github.com/mikulas/kubectl-repl)](https://goreportcard.com/report/github.com/mikulas/kubectl-repl)
[![Build Status](https://travis-ci.org/Mikulas/kubectl-repl.svg?branch=master)](https://travis-ci.org/Mikulas/kubectl-repl)

Wrap `kubectl` with namespace and variables.

[![asciicast](https://s3.eu-central-1.amazonaws.com/uploads.mangoweb.org/kubectl-repl-2.png)](https://asciinema.org/a/142536)


Installation
------------

Download latest release for your platform from https://github.com/Mikulas/kubectl-repl/releases.

Docker container is available at https://hub.docker.com/r/mikulas/kubectl-repl/. Requires volume mount into `/root/.kops`,
for example `-v ~/.kops:/root/.kops`.

Alternatively, download and build locally: see `Makefile` (`make build`). 


Usage
-----

`./kubectl-repl` first starts by asking you for namespace. You may enter any of the strings verbatim,
or any abbreviation that is closest. Additionally, you may use any of the variables REPL assigned (`$2`).

Then you are in the main REPL mode. You are presented with a prompt, into which you enter `kubectl` commands
(`kubectl -n $NS` prefix is implied).

The prompt can be exited with traditional *eof* or *sigint*, and an explicit `quit` or `exit` command. If a command
spawned long living process (such as `--follow`, `--watch` or `exec`) it will terminate the processes first and return
to repl.

I recommend using [rlwrap](https://github.com/hanslub42/rlwrap) in combination with `kubectl-repl`, such as
`rlwrap kubectl-repl`. This adds prompt history, search, buffering etc. Docker container distribution contains rlwrap.


Shell integration
-----------------

Instead of directly invoking `kubectl` with prompt as arguments, `/bin/sh -c` is used. This
allows for more complex usage usage as `grep` and redirects:

```console
# sentry get pods | grep app
+ kubectl -n sentry get pods | grep app
$1 	app-deployment-314667899-4r9c1       1/1       Running   0          22h
$2 	app-deployment-314667899-xr47k       1/1       Running   0          22h
# sentry get pods -o json > /tmp/pods.json
+ kubectl -n sentry get pods -o json > /tmp/pods.json
```


Raw shell invocation is also supported by prefixing the command with `;`. This should be intuitive as it works in
shell by default. Repl integration does not prefix the shell with `kubectl -n $NS` in this mode and trims the semicolon.
Repl variables are available as in all other prompt modes. 

```console
# kube-system get pods
+ kubectl -n kube-system get pods
   	NAME                        READY     STATUS    RESTARTS   AGE
$1 	kube-dns-3945342221-mwdh6   3/3       Running   0          9d
$2 	kube-dns-3945342221-x3fhn   3/3       Running   0          9d
# kube-system ; echo $(whoami) $2
+  echo $(whoami) kube-dns-3945342221-x3fhn
mikulas kube-dns-3945342221-x3fhn
```


Variables
---------

Prompts starting with `get` return their output prefixed with `$n`. You may use those variables anywhere to
automatically substitute for the value of the first column of the respective line. For example:
```console
# sentry get pods
+ kubectl -n sentry get pods
   	NAME                                 READY     STATUS    RESTARTS   AGE
$1 	app-deployment-314667899-4r9c1       1/1       Running   0          22h
$2 	app-deployment-314667899-xr47k       1/1       Running   0          22h
# sentry logs $2
+ kubectl -n sentry logs app-deployment-314667899-xr47k
```
The `$2` was substituted for `app-deployment-314667899-xr47k`.

Builtin variables have priority before shell variables from env, but both can be used: 

```console
$ env TYPE=pod rlwrap ./kubectl-repl -verbose
# default get $TYPE
+ kubectl -n default get $TYPE
```


Alternatives
------------

- https://github.com/c-bata/kube-prompt
