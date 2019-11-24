package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	clusterName string
	user        string
	pass        string
)

const (
	errorString            = "Unauthorized."
	pattern                = "/name"
	port                   = "8090"
	clusterNameUserVarName = "CLUSTER_NAME_U"
	clusterNamePassVarName = "CLUSTER_NAME_P"
)

func name(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s\n", clusterName)
}

func main() {
	if !envVarsValid() {
		fmt.Println("vars invalid")
		os.Exit(1)
	}
	user = os.Getenv(clusterNameUserVarName)
	pass = os.Getenv(clusterNamePassVarName)
	clusterName = getClusterName()
	if clusterName == "" {
		os.Exit(1)
	}
	http.HandleFunc(pattern, auth(name))
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func getClusterName() string {
	url := "http://metadata/computeMetadata/v1/instance/attributes/cluster-name"
	client := &http.Client{Timeout: time.Second * 10}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Metadata-Flavor", "Google")
	var resp *http.Response
	var err error
	if resp, err = client.Do(req); err != nil {

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

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
