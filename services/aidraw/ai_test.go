package aidraw

import (
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/utils/test"
	"testing"
)

func TestSexyAIDraw(t *testing.T) {

	if file.DataStorage.AiDraw.SexyAISession == "" {
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

func TestSexyAIOTP(t *testing.T) {
	res, err := SaiRequestOTP("abc@abc.com", false)
	if err != nil {
		t.Skip(err)
	}
	t.Logf("OTP: %+v", res)
}

func init() {
	test.InitTesting()
}
