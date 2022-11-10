package bilibili

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type parse struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

const parseJson = `
{
	"a": 1,
	"b": "str",
	"c": {
		"name": "Lam",
		"age": 15
	}
}
`

func TestParseMap(t *testing.T) {
	var gen interface{}
	if err := json.Unmarshal([]byte(parseJson), &gen); err != nil {
		t.Fatal(err)
	}
	if m, ok := gen.(map[string]interface{}); ok {
		if mm, okk := m["c"].(map[string]interface{}); okk {
			if b, err := json.Marshal(mm); err != nil {
				t.Fatal(err)
			} else {
				s := string(b)
				var p parse
				if err := json.Unmarshal([]byte(s), &p); err != nil {
					t.Fatal(err)
				} else {
					assert.Equal(t, p.Age, 15)
					assert.Equal(t, p.Name, "Lam")
				}
			}
		}
	}
}

type AAA struct {
	Text string
	A    int
}

var a = &AAA{
	Text: "abc",
	A:    123,
}

func TestPtrFunc(t *testing.T) {
	assert.Equal(t, "abc", a.Text)
	assignText(a)
	assert.Equal(t, "xyz", a.Text)
}

func assignText(a *AAA) {
	a.Text = "xyz"
}
