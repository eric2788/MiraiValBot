package set

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestStringSet(t *testing.T) {

	set := NewString()

	set.Add("hello")
	set.Add("world")
	set.Add("abc")
	set.Add("xyz")

	assert.Equal(t, set.Contains("xyz"), true)

	for v := range set.Iterator() {
		fmt.Println(v)
	}

	set.Delete("world")
	set.Delete("xyz")

	assert.Equal(t, set.Contains("xyz"), false)

	for v := range set.Iterator() {
		fmt.Println(v)
	}

}
