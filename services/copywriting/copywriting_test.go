package copywriting

import (
	"strings"
	"testing"
)

func TestCopyWriting(t *testing.T) {
	for _, line := range []string{"a %s", "b %s", "c"} {
		t.Logf(line, "hello world")
	}
}

func TestGetFabingList(t *testing.T) {
	list, _, err := GetFabingList()
	if err != nil {
		t.Skip(err)
	}
	t.Logf("%s", strings.Join(list, "\n"))
}

func TestGetFadianList(t *testing.T) {
	list, _, err := GetFadianList()
	if err != nil {
		t.Skip(err)
	}
	t.Logf("%s", strings.Join(list, "\n"))
}

func TestGetTiangouList(t *testing.T) {
	list, err := GetTianGouList()
	if err != nil {
		t.Skip(err)
	}
	t.Logf("%s", strings.Join(list, "\n"))
}
