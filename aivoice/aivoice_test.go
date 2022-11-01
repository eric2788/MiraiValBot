package aivoice

import (
	"os"
	"testing"
)

func TestGetGenshinVoice(t *testing.T) {
	b, err := GetGenshinVoice("別狗叫", "派蒙")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)
	err = os.MkdirAll("data", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("data/別狗叫.mp3", b, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
