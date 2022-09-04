package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

var (
	// Allows osExit to be stubbed during testing
	osExit = os.Exit
)

type ReleaseCleanup struct {
	action *githubactions.Action
	config Config
}

func NewReleaseCleanup() (*ReleaseCleanup, error) {
	action := githubactions.New()
	config, err := NewFromInputs(action)

	return &ReleaseCleanup{
		action: action,
		config: config,
	}, err
}

func (rc *ReleaseCleanup) Exit(err error) {
	if err != nil {
		rc.action.Errorf("%s", err.Error())
		osExit(1)
	} else {
		osExit(0)
	}
}

func (rc *ReleaseCleanup) Cleanup() error {
	filteredReleases, err := rc.getHelmReleases()
	if err != nil {
		return err
	}

	if len(filteredReleases) == 0 {
		rc.action.Infof("No releases found matching filters")
		return nil
	}

	rc.action.Infof(
		"The following releases matched the filters and will be deleted:\n%s",
		strings.Join(rc.prefixReleases("* ", filteredReleases), "\n"),
	)

	failedUninstalls := []string{}
	successfulUninstalls := []string{}
	for _, r := range filteredReleases {
		rc.header(fmt.Sprintf("Deleting release %s", r))
		args := rc.config.ToArgs(r)
		cmd := exec.Command("helm", args...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		if err := rc.runCommand(cmd); err != nil {
			rc.action.Errorf("Failed to delete release %s: %v", r, err)
			failedUninstalls = append(failedUninstalls, r)
		} else {
			rc.action.Infof("Successfully deleted release %s", r)
			successfulUninstalls = append(successfulUninstalls, r)
		}
	}

	rc.action.Infof("\n")
	rc.header("Summary")
	if len(successfulUninstalls) > 0 {
		rc.action.Infof("%d releases successfully deleted", len(successfulUninstalls))
	}

	if len(failedUninstalls) > 0 {
		return fmt.Errorf("%d releases failed to uninstall", len(failedUninstalls))
	}

	return nil
}

func (rc ReleaseCleanup) getHelmReleases() ([]string, error) {
	args := []string{"list", "--short"}

	if rc.config.KubeContext != "" {
		args = append(args, "--kube-context", rc.config.KubeContext)
	}

	if rc.config.Namespace != "" {
		args = append(args, "--namespace", rc.config.Namespace)
	}

	for _, f := range rc.config.Filters {
		args = append(args, f.Args()...)
	}

	out, err := rc.getCommandOutput(exec.Command("helm", args...))
	if err != nil {
		return []string{}, err
	}

	if out == "" {
		return []string{}, nil
	}

	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}

func (rc ReleaseCleanup) runCommand(cmd *exec.Cmd) error {
	rc.action.Infof("--> Running: %s", cmd.String())
	return cmd.Run()
}

func (rc ReleaseCleanup) getCommandOutput(cmd *exec.Cmd) (string, error) {
	rc.action.Infof("--> Running: %s", cmd.String())
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (rc ReleaseCleanup) prefixReleases(prefix string, releases []string) []string {
	prefixed := []string{}
	for _, r := range releases {
		prefixed = append(prefixed, fmt.Sprintf("%s%s", prefix, r))
	}
	return prefixed
}

func (rc ReleaseCleanup) header(title string) {
	titleLen := len(title)
	headerMask := strings.Repeat("-", titleLen+2)
	rc.action.Infof("%s\n %s\n%s", headerMask, title, headerMask)
}
