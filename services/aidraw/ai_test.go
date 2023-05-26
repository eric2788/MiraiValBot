package aidraw

import (
	"github.com/eric2788/MiraiValBot/utils/test"
	"testing"
)

func TestSexyAIDraw(t *testing.T) {

	sessionID = ""

	if sessionID == "" {
		t.Skip("sessionID is empty, skipped test")
	}

	res, err := sexyAIDraw(Payload{
		Prompt: "cat ears girl in the house",
		Model:  "real",
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("result image: %v", res.ImgUrl)

}

func init() {
	test.InitTesting()
}
