package main

import (
	"fmt"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	os.Setenv(envVarName("USER"), "TEST")
	os.Setenv(envVarName("PASS"), "TEST")
	if !envVarsValid() {
		t.Error("blah")
	}
}

func TestEnvVarName(t *testing.T) {
	if envVarName("TEST") != "CLUSTER_NAME_TEST" {
		t.Error("blah")
	}
}

var getEnvTests = []struct {
	key      string
	value    string
	fallback string
	expected string
	setEnv   bool
}{
	{"key", "value", "fallback-value", "fallback-value", false},
	{"key", "value", "fallback-value", "value", true},
}

func TestGetEnv(t *testing.T) {
	for _, getEnvTest := range getEnvTests {
		if getEnvTest.setEnv {
			os.Setenv(getEnvTest.key, getEnvTest.value)
		}
		if getEnv(getEnvTest.key, getEnvTest.fallback) != getEnvTest.expected {
			fmt.Println(getEnv(getEnvTest.key, getEnvTest.fallback))
			t.Error("blah")
		}
	}
}

var checkTests = []struct {
	suppliedUser string
	suppliedPass string
	user         string
	pass         string
	valid        bool
}{
	{"badUser", "badPass", "goodUser", "goodPass", false},
	{"goodUser", "goodPass", "goodUser", "goodPass", true},
}

func TestCheck(t *testing.T) {
	for _, checkTest := range checkTests {
		if check(checkTest.suppliedUser, checkTest.suppliedPass,
			checkTest.user, checkTest.pass) != checkTest.valid {
			t.Error("blah")
		}
	}

	// should pass as user/pass vars have been set
	// if !check(userPassDummy, userPassDummy) {
	// 	t.Error("blah")
	// }
}
