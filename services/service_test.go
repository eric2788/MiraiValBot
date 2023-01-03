package services

import (
	"strings"
	"testing"
)

func TestService(t *testing.T) {

	args := []string{"genshin", ""}

	tags := []string{""}

	if len(args) > 1 {
		tags = strings.Split(strings.Join(args, " "), ",")
	}

	t.Logf("tags: %v", tags)
}
