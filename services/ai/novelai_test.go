package ai

import (
	"testing"
)

func TestGetNovelAI8zywImage(t *testing.T) {

	url, err := GetNovelAI8zywImage(
		New8zywPayload(
			"1girl, best quality, masterpiece, cat ears, solo",
			WithoutR18,
		),
	)
	if err != nil {
		t.Skip(err)
	}

	t.Log(url)
}
