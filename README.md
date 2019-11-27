# cluster_name

This is a Go program destined for use within GKE clusters, that returns the name
of the cluster it's running in.

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
$ go build

$ ./cluster_name
```


If you choose the `go run` option, you'll either have to run it in the
background (add a `&` after the command), or open a new shell to test.

Now send a request to it:

```bash
$ curl -u blah:blah 'localhost:8090/cluster_name'

my_test_cluster_name
```