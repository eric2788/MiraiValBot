package datetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {

	before := time.Now().Unix()

	t.Logf("before: %d\n", before)

	<-time.After(time.Second * 5)

	after := time.Now().Unix()

	t.Logf("after: %d\n", after)

	assert.Equal(t, 5, Duration(before, after).Second())
}
