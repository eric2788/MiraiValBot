package file

import (
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var bilibili = &DataStorage.Bilibili

func TestLoadStorageReal(t *testing.T) {
	fmt.Println(len((*bilibili).HighLightedUsers.ToSlice()))
	DataStorage.Bilibili = &BilibiliSettings{
		HighLightedUsers: mapset.NewSet[int64](1, 2, 3),
	}
	fmt.Println(len((*bilibili).HighLightedUsers.ToSlice()), len(DataStorage.Bilibili.HighLightedUsers.ToSlice()))
	(*bilibili).HighLightedUsers.Add(9)
	fmt.Println(len((*bilibili).HighLightedUsers.ToSlice()), len(DataStorage.Bilibili.HighLightedUsers.ToSlice()))
}

var content = `
a: "hello world"
b: 1231
c: 
- "a"
- "b"
- "c"
`

func TestLoadYaml(t *testing.T) {
	_ = os.WriteFile("fake.yaml", []byte(content), 0644)
	var fakeYaml FakeYaml
	_ = loadYaml("fake.yaml", &fakeYaml)
	assert.Equal(t, "hello world", fakeYaml.A)
	assert.Equal(t, 1231, fakeYaml.B)
	assert.Equal(t, []string{"a", "b", "c"}, fakeYaml.C)
}

type boolJson struct {
	A bool `json:"a,string"`
	B bool `json:"b,string"`
}

func TestJsonParseBool(t *testing.T) {
	const test = `{"a": 1, "b": 0}`
	var m boolJson
	_ = json.Unmarshal([]byte(test), &m)
	t.Logf("A: %t, B: %t", m.A, m.B)
}

type FakeYaml struct {
	A string   `yaml:"a"`
	B int      `yaml:"b"`
	C []string `yaml:"c"`
}
