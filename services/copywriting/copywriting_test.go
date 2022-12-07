package copywriting

import "testing"

func TestCopyWriting(t *testing.T) {
	for _, line := range []string{"a %s", "b %s", "c"} {
		t.Logf(line, "hello world")
	}
}