package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/eric2788/common-utils/set"
	"github.com/stretchr/testify/assert"
)

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

type FakeYaml struct {
	A string   `yaml:"a"`
	B int      `yaml:"b"`
	C []string `yaml:"c"`
}
