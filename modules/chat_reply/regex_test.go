package chat_reply

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/eric2788/common-utils/array"
	//"github.com/stretchr/testify/assert"
)

func TestRegexp(t *testing.T) {
	face, err := regexp.Compile(`{face:(\d+)}`)
	if err != nil {
		t.Fatal(err)
	}

	content := "{face:8}好无聊哦,我都快睡着了{face:84}，聊的什么呀{face:1}蛤蛤"
	indexes := face.FindAllStringSubmatchIndex(content, -1)

	lastTo := 0
	for _, index := range indexes {
		//t.Log(index)
		from, to := index[0], index[1]

		if from > 0 {
			t.Logf("Append %s", content[lastTo:from])
		}

		fFrom, fTo := index[2], index[3]
		faceID, err := strconv.ParseInt(content[fFrom:fTo], 10, 32)
		if err == nil {
			t.Logf("Append Face %d", int32(faceID))
		} else {
			t.Logf("cannot append face: %v", err)
		}

		lastTo = to
	}

	if lastTo < len(content) {
		t.Logf("Append %s", content[lastTo:])
	}
}

func TestArrayAppend(t *testing.T) {
	a := []int{1, 2, 3}

	b := append([]int{}, a...)

	array.Remove(b, 1)

	t.Log(a)
	t.Log(b)
}
