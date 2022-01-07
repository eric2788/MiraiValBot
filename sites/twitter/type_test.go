package twitter

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func isA() bool {
	fmt.Println("is A")
	return false
}

func isB() bool {
	fmt.Println("is B")
	return false
}

func isC() bool {
	fmt.Println("is C")
	return false
}

func isD() bool {
	fmt.Println("is D")
	return true
}

func TestSwitch(t *testing.T) {
	switch {
	case isB():
		fmt.Println("B")
		return
	case isA():
		fmt.Println("A")
		return
	case isC():
		fmt.Println("C")
		return
	case isD():
		fmt.Println("D")
		return
	default:
		fmt.Println("Nil")
		return
	}
}

type TestPtr struct {
	A *struct {
		B int16
		C *string
	}
}

var testPtrStr = `
{
	"A": {
		"B": 12345
    }
}
`

var testPtrStr2 = `{}`

func TestPtrStruct(t *testing.T) {
	var testPtr = &TestPtr{}
	if err := json.Unmarshal([]byte(testPtrStr), testPtr); err != nil {
		t.Fatal(err)
	}

	fmt.Println(*testPtr.A)

	assert.Nil(t, testPtr.A.C)

	testPtr = &TestPtr{}
	if err := json.Unmarshal([]byte(testPtrStr2), testPtr); err != nil {
		t.Fatal(err)
	}

	fmt.Println(*testPtr)

	assert.Nil(t, testPtr.A)
}
