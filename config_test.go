package main

import (
	"os"
	"strings"
	"testing"

	"github.com/c2fo/testify/assert"
	"github.com/sethvargo/go-githubactions"
)

func withEnv(t *testing.T, env map[string]string, f func()) {
	// Save the current environment
	preTest := map[string]string{}
	for _, e := range os.Environ() {
		kv := strings.Split(e, "=")
		preTest[kv[0]] = kv[1]
	}

	for k, v := range env {
		if err := os.Setenv(k, v); err != nil {
			t.Fatalf("failed to set env var %s: %s", k, err)
		}
	}

	f()

	for k := range env {
		if _, ok := preTest[k]; !ok {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, preTest[k])
		}
	}
}

func TestValuesFiles(t *testing.T) {
	testEnv := map[string]string{
		"INPUT_NAMESPACE":                "mynamespace",
		"INPUT_KUBE-CONTEXT":             "mycontext",
		"INPUT_LAST-MODIFIED-OLDER-THAN": "1h",
	}

	withEnv(t, testEnv, func() {
		config, err := NewFromInputs(githubactions.New())
		assert.NoError(t, err)

		assert.Equal(t, "mynamespace", config.Namespace)
		assert.Equal(t, "mycontext", config.KubeContext)
		assert.NotEmpty(t, config.Filters)
	})
}
