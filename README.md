Kubectl REPL
============

[![Go Report Card](https://goreportcard.com/badge/github.com/mikulas/kubectl-repl)](https://goreportcard.com/report/github.com/mikulas/kubectl-repl)
[![Build Status](https://travis-ci.org/Mikulas/kubectl-repl.svg?branch=master)](https://travis-ci.org/Mikulas/kubectl-repl)

Wrap `kubectl` with namespace and variables.

[![asciicast](https://asciinema.org/a/142536.png)](https://asciinema.org/a/142536)


Installation
------------

Download latest release for your platform from https://github.com/Mikulas/kubectl-repl/releases.

Alternatively, download and build locally: see `Makefile` (`make build`). 


Usage
-----

`./kubectl-repl` first starts by asking you for namespace. You may enter any of the strings verbatim,
or any abbreviation that is closest. Additionally, you may use any of the variables REPL assigned (`$2`).

Then you are in the main REPL mode. You are presented with a prompt, into which you enter `kubectl` commands
(`kubectl -n $NS` prefix is implied).

The prompt can be exited with traditional *eof* or *sigint*, and an explicit `quit` or `exit` command.

I recommend using [rlwrap](https://github.com/hanslub42/rlwrap) in combination with `kubectl-repl`, such as
`rlwrap kubectl-repl`. This adds prompt history, search, buffering etc.


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
