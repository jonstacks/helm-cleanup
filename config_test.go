package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sethvargo/go-githubactions"
	"github.com/stretchr/testify/assert"
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

func TestConfigNewFromInputs(t *testing.T) {
	testEnv := map[string]string{
		"INPUT_NAMESPACE":                "mynamespace",
		"INPUT_KUBE-CONTEXT":             "mycontext",
		"INPUT_LAST-MODIFIED-OLDER-THAN": "168h", // 7 days
	}

	withEnv(t, testEnv, func() {
		config, err := NewFromInputs(githubactions.New())
		assert.NoError(t, err)

		assert.Equal(t, "mynamespace", config.Namespace)
		assert.Equal(t, "mycontext", config.KubeContext)
		assert.Nil(t, config.Timeout)
		assert.NotEmpty(t, config.Filters)
	})

	testEnv = map[string]string{
		"INPUT_LAST-MODIFIED-OLDER-THAN": "24h",
		"INPUT_RELEASE-NAME-REGEX":       "myapp-.*",
		"INPUT_DRY-RUN":                  "true",
		"INPUT_NO-HOOKS":                 "true",
		"INPUT_KEEP-HISTORY":             "false",
		"INPUT_WAIT":                     "true",
		"INPUT_TIMEOUT":                  "10m",
	}

	withEnv(t, testEnv, func() {
		config, err := NewFromInputs(githubactions.New())
		assert.NoError(t, err)

		assert.Equal(t, true, config.DryRun)
		assert.Equal(t, true, config.NoHooks)
		assert.Equal(t, false, config.KeepHistory)
		assert.Equal(t, true, config.Wait)
		assert.Equal(t, 10*time.Minute, *config.Timeout)
	})

	// Test invalid timeout with unsupported duration format
	withEnv(t, map[string]string{"INPUT_TIMEOUT": "7d"}, func() {
		_, err := NewFromInputs(githubactions.New())
		assert.Error(t, err, "Unable to parse timeout as a duration. See https://pkg.go.dev/time#ParseDuration for valid duration formats.")
	})

	// Test invalid timeout with bad durations
	withEnv(t, map[string]string{"INPUT_TIMEOUT": "aaaa"}, func() {
		_, err := NewFromInputs(githubactions.New())
		assert.Error(t, err, "Unable to parse timeout as a duration. See https://pkg.go.dev/time#ParseDuration for valid duration formats.")
	})

	// Test invalid last-modified-older-than
	withEnv(t, map[string]string{"INPUT_LAST-MODIFIED-OLDER-THAN": "3d"}, func() {
		_, err := NewFromInputs(githubactions.New())
		assert.Error(t, err, "Unable to parse last-modified-older-than as a duration. See https://pkg.go.dev/time#ParseDuration for valid duration formats.")
	})

	// Test no Filters
	withEnv(t, map[string]string{"INPUT_DRY-RUN": "true"}, func() {
		config, err := NewFromInputs(githubactions.New())
		assert.Empty(t, config.Filters)
		assert.Error(t, err, "One of release-name-filter or last-modified-older-than are required.")
	})

}

func TestConfigToArgs(t *testing.T) {
	now := time.Unix(1615392800, 0)
	timeout := 10 * time.Minute

	config := Config{
		Namespace: "mynamespace",
		Timeout:   &timeout,
		Debug:     true,
		Wait:      true,
		Filters: []Filter{
			ModifiedAtLessThanFilter{
				Now:      now,
				Lookback: 24 * time.Hour,
			},
		},
	}

	assert.Equal(t,
		[]string{"uninstall", "my-release", "--debug", "--namespace", "mynamespace", "--timeout", "10m0s", "--wait"},
		config.ToArgs("my-release"),
	)
}
