package aivoice

import (
	"os"
	"testing"
)

func TestGetGenshinVoice(t *testing.T) {
	b, err := GetGenshinVoice("別狗叫", "派蒙")
	if err != nil {
		t.Skip(err)
	}
	if len(b) == 0 {
		t.Fatal("empty voice")
	}
	err = os.MkdirAll("data", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("data/別狗叫.silk", b, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
