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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	defaultPattern = "/name"
	defaultPort    = "8090"
	envVarPrefix   = "CLUSTER_NAME"
	userVarName    = "USER"
	passVarName    = "PASS"
	patternVarName = "PATT"
	portVarName    = "PORT"
	metaURL        = "http://metadata/computeMetadata/v1/instance/attributes/cluster-name"
)

func name(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s\n", clusterName)
}

func main() {
	if !envVarsValid() {
		logger.Fatal("Invalid env vars")
	}
	user = os.Getenv(envVarName(userVarName))
	pass = os.Getenv(envVarName(passVarName))
	var err error
	if clusterName, err = getClusterName(); err != nil {
		logger.Fatal(err)
	}
	pattern := getEnv(envVarName(patternVarName), defaultPattern)
	http.HandleFunc(pattern, auth(name))
	port := getEnv(envVarName(patternVarName), defaultPort)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func getClusterName() (name string, err error) {
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
	return
}

func envVarsValid() (valid bool) {
	return len(os.Getenv(envVarName(userVarName))) > 0 &&
		len(os.Getenv(envVarName(passVarName))) > 0
}

func envVarName(suffix string) string {
	return fmt.Sprintf("%s_%s", envVarPrefix, suffix)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

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

func check(basicUser, basicPass string) bool {
	return basicUser == user && basicPass == pass
}
