package cmd

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestLongTu(t *testing.T) {
	backup := "https://phqghume.github.io/img/"
	rand.Seed(time.Now().UnixMicro())
	random := rand.Intn(58) + 1
	ext := ".jpg"
	if random > 48 {
		ext = ".gif"
	}
	imgLink := fmt.Sprintf("%slong%%20(%d)%s", backup, random, ext)
	t.Log(imgLink)
}

func TestPanic(t *testing.T) {
	a := []string{"a"}
	t.Log(a[1:]) // no panic!
}
