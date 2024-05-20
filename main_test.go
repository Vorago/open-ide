package main

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func Test_projectRegex(t *testing.T) {
	r := projectRegexp("open-ide")
	match, _ := regexp.MatchString(r, "open-ide – main.go")
	assert.True(t, match)

	match, _ = regexp.MatchString(r, "open-ide")
	assert.True(t, match)

	match, _ = regexp.MatchString(r, "another-project – open-ide")
	assert.False(t, match)
}
