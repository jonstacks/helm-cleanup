package main

import (
	"fmt"
	"strconv"
	"time"
)

type Filter interface {
	Args() []string
}

func timeToUnixString(t time.Time) string {
	return strconv.Itoa(int(t.Unix()))
}

type ModifiedAtLessThanFilter struct {
	Now      time.Time
	Lookback time.Duration
}

func (f ModifiedAtLessThanFilter) Args() []string {
	modifiedAt := timeToUnixString(f.Now.Add(-f.Lookback))
	return []string{"--selector", fmt.Sprintf("modifiedAt<%s", modifiedAt)}
}

type ReleaseNameFilter struct {
	FilterString string
}

func (f ReleaseNameFilter) Args() []string {
	return []string{"--filter", f.FilterString}
}
