package main

import (
	"errors"
	"time"

	"github.com/sethvargo/go-githubactions"
)

func parseBool(action *githubactions.Action, name string) bool {
	return action.GetInput(name) == "true"
}

type Config struct {
	// Filters
	Filters []Filter

	// Optional
	Debug       bool
	Description string
	DryRun      bool
	KeepHistory bool
	KubeContext string
	Namespace   string
	NoHooks     bool
	Timeout     *time.Duration
	Wait        bool
}

func NewFromInputs(action *githubactions.Action) (Config, error) {
	c := &Config{}

	c.Debug = parseBool(action, "debug")
	c.Description = action.GetInput("description")
	c.DryRun = parseBool(action, "dry-run")
	c.KeepHistory = parseBool(action, "keep-history")
	c.KubeContext = action.GetInput("kube-context")
	c.Namespace = action.GetInput("namespace")
	c.NoHooks = parseBool(action, "no-hooks")
	c.Wait = parseBool(action, "wait")

	timeout := action.GetInput("timeout")
	if timeout != "" {
		d, err := time.ParseDuration(timeout)
		if err != nil {
			return *c, errors.New("Unable to parse timeout as a duration. See https://pkg.go.dev/time#ParseDuration for valid duration formats.")
		}
		c.Timeout = &d
	}

	releaseNameFilter := action.GetInput("release-name-filter")
	if releaseNameFilter != "" {
		c.Filters = append(c.Filters, ReleaseNameFilter{FilterString: releaseNameFilter})
	}

	lastModifiedOlderThan := action.GetInput("last-modified-older-than")
	if lastModifiedOlderThan != "" {
		d, err := time.ParseDuration(lastModifiedOlderThan)
		if err != nil {
			return *c, errors.New("Unable to parse last-modified-older-than as a duration. See https://pkg.go.dev/time#ParseDuration for valid duration formats.")
		}
		c.Filters = append(c.Filters, ModifiedAtLessThanFilter{
			Now:      time.Now(),
			Lookback: d,
		})
	}

	if len(c.Filters) == 0 {
		return *c, errors.New("One of release-name-filter or last-modified-older-than are required.")
	}

	return *c, nil
}

func (c Config) ToArgs(releaseName string) []string {
	args := []string{"uninstall", releaseName}

	if c.Debug {
		args = append(args, "--debug")
	}

	if c.Description != "" {
		args = append(args, "--description", c.Description)
	}

	if c.DryRun {
		args = append(args, "--dry-run")
	}

	if c.KeepHistory {
		args = append(args, "--keep-history")
	}

	if c.KubeContext != "" {
		args = append(args, "--kube-context", c.KubeContext)
	}

	if c.Namespace != "" {
		args = append(args, "--namespace", c.Namespace)
	}

	if c.NoHooks {
		args = append(args, "--no-hooks")
	}

	if c.Timeout != nil {
		args = append(args, "--timeout", c.Timeout.String())
	}

	if c.Wait {
		args = append(args, "--wait")
	}

	return args
}
