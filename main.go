// Copyright 2019 OVO Technology
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/subtle"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ovotech/cluster_name/pkg/log"
)

var (
	clusterName string
	user        string
	pass        string
	logger      = log.StdoutLogger().Sugar()
)

const (
	errorString    = "Unauthorized."
	defaultPattern = "/cluster_name"
	defaultPort    = "8090"
	defaultLocal   = "false"
	envVarPrefix   = "CLUSTER_NAME"
	userVarName    = "USER"
	passVarName    = "PASS"
	patternVarName = "PATT"
	portVarName    = "PORT"
	localVarName   = "LOCAL"
	metaURL        = "http://metadata/computeMetadata/v1/instance/attributes/cluster-name"
)

func main() {
	if !envVarsValid() {
		logger.Fatal("Invalid env vars")
	}
	user = os.Getenv(envVarName(userVarName))
	pass = os.Getenv(envVarName(passVarName))
	var err error
	if clusterName, err = getClusterName(getEnv(envVarName(localVarName), defaultLocal)); err != nil {
		logger.Fatal(err)
	}
	pattern := getEnv(envVarName(patternVarName), defaultPattern)
	http.HandleFunc(pattern, auth(name))
	port := getEnv(envVarName(patternVarName), defaultPort)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

// name writes to the provided writer
func name(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, clusterName)
}

// auth checks basic auth is successful before calling the provided func
func auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		basicUser, basicPass, _ := r.BasicAuth()
		if !check(basicUser, basicPass) {
			logger.Warnf("Failed auth attempt from %s", r.RemoteAddr)
			http.Error(w, errorString, http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}

// getClusterName gets the cluster name from the Kubernetes API
func getClusterName(localString string) (name string, err error) {
	var local bool
	local, err = strconv.ParseBool(localString)
	if err != nil {
		return
	}
	if local {
		name = "my_test_cluster_name"
	} else {
		url := metaURL
		client := &http.Client{Timeout: time.Second * 10}
		var req *http.Request
		if req, err = http.NewRequest("GET", url, nil); err != nil {
			return
		}
		req.Header.Set("Metadata-Flavor", "Google")
		var resp *http.Response
		if resp, err = client.Do(req); err != nil {
			return
		}
		defer resp.Body.Close()
		var body []byte
		if body, err = ioutil.ReadAll(resp.Body); err != nil {
			return
		}
		name = string(body)
	}
	return
}

// envVarsValid returns true if all required env vars are set
func envVarsValid() (valid bool) {
	return len(os.Getenv(envVarName(userVarName))) > 0 &&
		len(os.Getenv(envVarName(passVarName))) > 0
}

// envVarName returns an env var name created from a prefix and suffix
func envVarName(suffix string) string {
	return fmt.Sprintf("%s_%s", envVarPrefix, suffix)
}

// getEnv returns the value of the specified env var, or the fallback string if
// it hasn't been set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// check returns true if the provided user and pass strings are equal to
// their configured counterparts
func check(basicUser, basicPass string) bool {
	return subtle.ConstantTimeCompare([]byte(basicUser), []byte(user)) == 1 &&
		subtle.ConstantTimeCompare([]byte(basicPass), []byte(pass)) == 1
}
