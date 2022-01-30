package file

import (
	"encoding/json"
	"fmt"
	"github.com/eric2788/common-utils/set"
	"github.com/stretchr/testify/assert"
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
	E *fakeNested
}

type fakeNested struct {
	E string
}

var defaultFakeStorageData = &fakeStorageData{
	A: []int16{1, 2},
	B: []int16{3, 4},
	C: []string{},
	D: []string{"Default", "Data"},
	E: &fakeNested{
		E: "hiawhdiawhiw",
	},
}

var e = &defaultFakeStorageData.E

func TestStorageSwitch(t *testing.T) {
	fmt.Printf("before: %v\n", (*e).E)
	defaultFakeStorageData.E = &fakeNested{
		E: "HELLO",
	}
	fmt.Printf("after: %v\n", (*e).E)
	assert.Equal(t, defaultFakeStorageData.E.E, (*e).E)
}

func TestLoadStorage(t *testing.T) {
	var fakeStorageData = defaultFakeStorageData
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

var bilibili = &DataStorage.Bilibili

func TestLoadStorageReal(t *testing.T) {
	fmt.Println((*bilibili).HighLightedUsers.Size())
	DataStorage.Bilibili = &BilibiliSettings{
		HighLightedUsers: set.FromInt64Arr([]int64{1, 2, 3}),
	}
	fmt.Println((*bilibili).HighLightedUsers.Size(), DataStorage.Bilibili.HighLightedUsers.Size())
	(*bilibili).HighLightedUsers.Add(9)
	fmt.Println((*bilibili).HighLightedUsers.Size(), DataStorage.Bilibili.HighLightedUsers.Size())
}
