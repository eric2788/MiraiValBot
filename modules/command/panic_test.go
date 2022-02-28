package command

import (
	"fmt"
	"testing"
)

func ATestPanic(t *testing.T) {
	invoke("A")
	invoke("B")
	invoke("C")
}

func invoke(t string) {
	defer func() {
		fmt.Println(recover())
	}()
	if t == "B" {
		panic("B!")
	}
}
