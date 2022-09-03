package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

func GetHelmReleases(context, namespace string, filters []Filter) ([]string, error) {
	args := []string{"list", "--short"}

	for _, f := range filters {
		args = append(args, f.Args()...)
	}

	cmd := exec.Command("helm", args...)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(string(out), "\n"), nil
}

func main() {
	action := githubactions.New()
	config, err := NewFromInputs(action)
	if err != nil {
		action.Fatalf("%v", err)
	}

	filteredReleases, err := GetHelmReleases(config.KubeContext, config.Namespace, config.Filters)
	if err != nil {
		action.Fatalf("%v", err)
	}

	if len(filteredReleases) == 0 {
		action.Infof("No releases found matching filters")
		os.Exit(0)
	}

	failedUninstalls := []string{}
	successfulUninstalls := []string{}
	for _, r := range filteredReleases {
		action.Infof("Deleting release %s", r)
		args := config.ToArgs(r)
		cmd := exec.Command("helm", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			action.Errorf("Failed to delete release %s: %v", r, err)
			failedUninstalls = append(failedUninstalls, r)
		} else {
			successfulUninstalls = append(successfulUninstalls, r)
		}
	}

	if len(successfulUninstalls) > 0 {
		action.Infof(
			"Successfully deleted the following %d releases:\n%s",
			len(successfulUninstalls),
			strings.Join(successfulUninstalls, "\n  "),
		)
	}

	if len(failedUninstalls) > 0 {
		action.Fatalf(
			"The following %d released failed to uninstall: \n%s",
			len(failedUninstalls),
			strings.Join(failedUninstalls, "\n"),
		)
	}
}
