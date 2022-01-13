package bilibili

import (
	"encoding/json"
	"fmt"
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
					fmt.Printf("%+v", p)
				}
			}
		}
	}
}

func TestMarshalMap(t *testing.T) {
	a := map[string]interface{}{
		"a": 1,
		"b": "str",
	}
	b, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}
