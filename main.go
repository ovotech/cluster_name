package main

import (
	"fmt"
	"net/http"
	"os"
)

var clusterName string
var user string
var pass string
var errorString = "Unauthorized."
var pattern = "/name"
var port = "8090"
var clusterNameUserVarName = "CLUSTER_NAME_U"
var clusterNamePassVarName = "CLUSTER_NAME_P"

func name(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s\n", clusterName)
}

func main() {
	if !envVarsValid() {
		os.Exit(1)
	}
	//TODO: get the actual cluster name instead of this
	clusterName = "burg"
	if clusterName == "" {
		os.Exit(1)
	}
	http.HandleFunc(pattern, auth(name))
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

// curl http://metadata/computeMetadata/v1/instance/attributes/cluster-name -H "Metadata-Flavor: Google"

func envVarsValid() (valid bool) {
	return len(os.Getenv(clusterNameUserVarName)) > 0 &&
		len(os.Getenv(clusterNamePassVarName)) > 0
}

func auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		basicUser, basicPass, _ := r.BasicAuth()
		if !check(basicUser, basicPass) {
			http.Error(w, errorString, http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}

func check(basicUser, basicPass string) bool {
	return basicUser == user && basicPass == pass
}
