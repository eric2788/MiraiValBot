package file

import (
	"encoding/json"
	"fmt"
	"testing"
)

const fakeStorage = `
{
	"a": [1, 2, 3],
	"b": [4, 5, 6],
	"c": ["a", "b", "c"]
}
`

type FakeStorageData struct {
	A []int16  `json:"a"`
	B []int16  `json:"b"`
	C []string `json:"c"`
	D []string `json:"d"`
}

var defaultFakeStorageData = FakeStorageData{
	A: []int16{},
	B: []int16{},
	C: []string{},
	D: []string{"Default", "Data"},
}

func TestLoadStorage(t *testing.T) {
	var fakeStorageData = &defaultFakeStorageData
	if err := json.Unmarshal([]byte(fakeStorage), fakeStorageData); err != nil {
		t.Fatal(err)
	}
	fmt.Println(*fakeStorageData)
	fmt.Println(defaultFakeStorageData)

}
