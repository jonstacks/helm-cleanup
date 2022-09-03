package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/sethvargo/go-githubactions"
)

func parseBool(action *githubactions.Action, name string) bool {
	return action.GetInput(name) == "true"
}

func parseStringSlice(action *githubactions.Action, name, delimiter string) []string {
	return strings.Split(action.GetInput(name), delimiter)
}

type Filter interface {
	Args() []string
}

type LabelFilter struct {
	Label    string
	Value    string
	Operator string
}

func (f LabelFilter) Args() []string {
	return []string{"--selector", fmt.Sprintf("%s%s%s", f.Label, f.Operator, f.Value)}
}

type ReleaseNameFilter struct {
	Regexp *regexp.Regexp
}

func (f ReleaseNameFilter) Args() []string {
	return []string{"--filter", f.Regexp.String()}
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

	timeout := action.GetInput("timeout")
	if timeout != "" {
		d, err := time.ParseDuration(timeout)
		if err != nil {
			return *c, errors.New("Unable to parse timeout")
		}
		c.Timeout = &d
	}

	releaseNameFilter := action.GetInput("release-name-filter")
	if releaseNameFilter != "" {
		re, err := regexp.Compile(releaseNameFilter)
		if err != nil {
			return *c, errors.New("Unable to parse release-name-filter as a regular expression")
		}
		c.Filters = append(c.Filters, ReleaseNameFilter{Regexp: re})
	}

	lastModifiedOlderThan := action.GetInput("last-modified-older-than")
	if lastModifiedOlderThan != "" {
		d, err := time.ParseDuration(lastModifiedOlderThan)
		if err != nil {
			return *c, errors.New("Unable to parse last-modified-older-than as a duration")
		}
		c.Filters = append(c.Filters, LabelFilter{
			Label:    "modifiedAt",
			Operator: "<",
			Value:    fmt.Sprintf("%d", time.Now().Unix()-int64(d.Seconds())),
		})
	}

	if len(c.Filters) == 0 {
		return *c, errors.New("One of release-name-filter or last-modified-older-than are required")
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
