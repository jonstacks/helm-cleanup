//go:build integration
// +build integration

package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var repoAddStable = exec.Command("helm", "repo", "add", "stable", "https://charts.helm.sh/stable")
var repoAddBitnami = exec.Command("helm", "repo", "add", "bitnami", "https://charts.bitnami.com/bitnami")
var repoUpdate = exec.Command("helm", "repo", "update")
var createTestRelease = func(releaseName, namespace string) *exec.Cmd {
	return exec.Command(
		"helm", "upgrade",
		"--install",
		"--debug",
		"--create-namespace",
		"--namespace", namespace,
		"--wait",
		"--timeout", "3m",
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
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Run()
		if err != nil {
			assert.FailNow(t, err.Error())
		}
	}

	testRelease := func(i int) string { return fmt.Sprintf("test-release-%d", i) }

	runCommand(repoAddStable)
	runCommand(repoAddBitnami)
	runCommand(repoUpdate)

	for i := 0; i < 5; i++ {
		runCommand(createTestRelease(testRelease(i), "helm-cleanup-1"))
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

		for i := 0; i < 3; i++ {
			assert.NotContains(t, releases, testRelease(i))
		}

		for i := 3; i < 5; i++ {
			assert.Contains(t, releases, testRelease(i))
		}
	})

	runCommand(cleanupNamespace("helm-cleanup-1"))
}
