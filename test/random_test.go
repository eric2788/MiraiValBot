package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRandom(t *testing.T) {
	rand.Seed(time.Now().UnixMicro())
	for i := 0; i < 200; i++ {
		n := rand.Intn(2)
		fmt.Println(n)
		if n == 2 {
			t.Fatal("2")
		}
	}
}
