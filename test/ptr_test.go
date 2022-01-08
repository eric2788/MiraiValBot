package test

import (
	"context"
	"fmt"
	"testing"
)

type A struct {
	B *B
}

type B struct {
	C *C
}

type C struct {
	Arr []string
}

var a = &A{
	B: &B{
		C: &C{
			Arr: []string{},
		},
	},
}

var c = a.B.C
var c2 = a.B.C

func TestPtr(t *testing.T) {
	c.Arr = append(c.Arr, "hello")
	fmt.Println(a.B.C.Arr)
	c2.Arr = append(c2.Arr, "world")
	fmt.Println(a.B.C.Arr)
}

func TestCtxPtr(t *testing.T) {
	var s *context.Context

	fmt.Println(s == nil)

	assignCtx(s)

	fmt.Println(s == nil)
}

func assignCtx(s *context.Context) {
	ctx := context.Background()
	s = &ctx
}
