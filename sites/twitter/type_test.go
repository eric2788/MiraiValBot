package twitter

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
