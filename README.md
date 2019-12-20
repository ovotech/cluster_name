# cluster_name
[![CircleCI](https://circleci.com/gh/ovotech/cluster_name/tree/master.svg?style=svg)](https://circleci.com/gh/ovotech/cluster_name/tree/master)

This is a Go program designed for use within GKE clusters, that returns the name of the cluster it's running in.

The cluster name is obtained by making a call to the [Compute Engine API](https://cloud.google.com/compute/docs/storing-retrieving-metadata):

```
http://metadata/computeMetadata/v1/instance/attributes/cluster-name
```

## Install

- Build from source (see below)
- Binaries on [GitHub releases](https://github.com/ovotech/cluster_name/releases)
- [Docker image](https://hub.docker.com/r/ovotech/cluster_name)

## Running Locally

Set the required env vars:

```bash
$ export CLUSTER_NAME_USER=foo CLUSTER_NAME_PASS=bar CLUSTER_NAME_LOCAL=true
```

Check out the project, and in its root:

```bash
$ go run main.go
```

OR, build a binary, and run it

```bash
$ export GO111MODULE=on 

$ go build

$ ./cluster_name
```

You'll either have to run it in the background (add a `&` after the command), 
or open a new shell to test.

Now send a request to it:

```bash
$ curl -u foo:bar 'localhost:8091/cluster_name'
my_test_cluster_name
```

## Configuration

| Env var       | Description           | Required?  | Default |
|:-------------:|:-------------:|:-----:|:-----:|
| CLUSTER_NAME_LOCAL  | Boolean indicating app is being run locally | n | false |
| CLUSTER_NAME_PASS  | Password string | y | |
| CLUSTER_NAME_PATH  | Path that the app should respond to | n | "/cluster_name" |
| CLUSTER_NAME_PORT  | Port that the app should listen to | n | "8091" |
| CLUSTER_NAME_USER  | Username string | y | |

## Healthz

There's also a healthz endpoint that simply returns a 200 status once the application has started:

```bash
$ curl -s -o /dev/null -w "%{http_code}" localhost:8091/healthz
200
```


## Auth

Currently only basic auth is supported.


## TLS

TLS is not supported.