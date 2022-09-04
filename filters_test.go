package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestModifiedAtLessThanFilterImplementsFilter(t *testing.T) {
	assert.Implements(t, (*Filter)(nil), new(ModifiedAtLessThanFilter))
}

func TestModifiedAtLessThanFilter(t *testing.T) {
	now := time.Unix(1615392800, 0)

	filter := ModifiedAtLessThanFilter{
		Now:      now,
		Lookback: 7 * 24 * time.Hour,
	}

	assert.Equal(t, []string{"--selector", "modifiedAt<1614788000"}, filter.Args())
}

func TestReleaseNameFilterImplementsFilter(t *testing.T) {
	assert.Implements(t, (*Filter)(nil), new(ReleaseNameFilter))
}

func TestReleaseNameFilter(t *testing.T) {
	filter := ReleaseNameFilter{FilterString: `^myapp-.*`}
	assert.Equal(t, []string{"--filter", "^myapp-.*"}, filter.Args())
}
