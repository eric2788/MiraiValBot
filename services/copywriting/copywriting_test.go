package copywriting

import (
	"math/rand"
	"strings"
	"testing"
	"time"
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

func TestGetCPList(t *testing.T) {
	list, _, _, err := GetCPList()
	if err != nil {
		t.Skip(err)
	}
	t.Logf("%s", strings.Join(list, "\n"))
}

func TestGetCrazyThursdayList(t *testing.T) {
	list, err := GetCrazyThursdayList()
	if err != nil {
		t.Skip(err)
	}
	t.Logf("%s", strings.Join(list, "\n"))
}

func TestGetRanran(t *testing.T) {
	list, err := GetRanranList()
	if err != nil {
		t.Skip(err)
	}
	for _, as := range list {
		t.Logf(strings.ReplaceAll(as.Text, as.Person, "夏诺雅"))
	}
}

func TestRandom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	t.Log(rand.Intn(100))
	t.Log(rand.Intn(100))
	//rand.Seed(time.Now().UnixNano())
	t.Log(rand.Intn(100))
	t.Log(rand.Intn(100))
	t.Log(rand.Intn(100))
	t.Log(rand.Intn(100))
	t.Log(rand.Intn(100))
}
