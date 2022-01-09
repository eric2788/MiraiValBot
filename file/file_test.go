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

type fakeStorageData struct {
	A []int16  `json:"a"`
	B []int16  `json:"b"`
	C []string `json:"c"`
	D []string `json:"d"`
	E fakeNested
}

type fakeNested struct {
	E string
}

var defaultFakeStorageData = fakeStorageData{
	A: []int16{1, 2},
	B: []int16{3, 4},
	C: []string{},
	D: []string{"Default", "Data"},
	E: fakeNested{
		E: "hiawhdiawhiw",
	},
}

func TestLoadStorage(t *testing.T) {
	var fakeStorageData = &defaultFakeStorageData
	if err := json.Unmarshal([]byte(fakeStorage), fakeStorageData); err != nil {
		t.Fatal(err)
	}
	fmt.Println(*fakeStorageData)
	fmt.Println(defaultFakeStorageData)
}

func TestSaveStorage(t *testing.T) {
	b, err := json.Marshal(defaultFakeStorageData)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}
