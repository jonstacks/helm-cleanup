//go:build integration
// +build integration

package main

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var repoAdd = exec.Command("helm", "repo", "add", "stable", "https://charts.helm.sh/stable")
var repoUpdate = exec.Command("helm", "repo", "update")
var createTestRelease = func(releaseName, namespace string) *exec.Cmd {
	return exec.Command(
		"helm", "upgrade",
		"--install",
		"--atomic",
		"--create-namespace",
		"--namespace", namespace,
		"--wait",
		"--timeout", "5m",
		"--set", "service.type=ClusterIP",
		releaseName, "bitnami/nginx",
	)
}
var cleanupNamespace = func(namespace string) *exec.Cmd {
	return exec.Command("kubectl", "delete", "namespace", namespace)
}

var getReleases = func(namespace string) (string, error) {
	b, err := exec.Command("helm", "list", "--namespace", namespace).CombinedOutput()
	return string(b), err
}

func TestCleansUpFilteredReleases(t *testing.T) {
	runCommand := func(c *exec.Cmd) {
		assert.NoError(t, c.Run())
	}

	runCommand(repoAdd)
	runCommand(repoUpdate)

	for i := 0; i < 5; i++ {
		releaseName := fmt.Sprintf("test-release-%d", i)
		runCommand(createTestRelease(releaseName, "helm-cleanup-1"))
	}

	testEnv := map[string]string{
		"INPUT_NAMESPACE":           "helm-cleanup-1",
		"INPUT_RELEASE-NAME-FILTER": `test-release-[0-2]`,
		"INPUT_DEBUG":               "true",
	}

	withEnv(t, testEnv, func() {
		c, err := NewReleaseCleanup()
		assert.Nil(t, err)

		assert.Nil(t, c.Cleanup())

		releases, err := getReleases("helm-cleanup-1")
		assert.NoError(t, err)
		assert.Contains(t, releases, "test-release-3")
		assert.NotEmpty(t, releases, "test-release-4")
	})

	runCommand(cleanupNamespace("helm-cleanup-1"))
}
