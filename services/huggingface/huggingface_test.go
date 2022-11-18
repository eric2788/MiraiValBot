package huggingface

import (
	"os"
	"testing"
)

func TestWaifuDiffusier(t *testing.T) {

	if os.Getenv("HUGGING_FACE_TOKEN") == "" {
		t.Skip("no token set")
	}

	b, err := GetResultImage("Nilaier/Waifu-Diffusers",
		NewParam(
			Input("group of girls, golden details, gems, fluffy hair, silver hair straight high ponytail+long hair, curly blonde hair, {ginger hair}, purple hair low ponytail curly, clear details, night time, fireworks, casual clothes, raytracing, foreground focus, blurred background, intricate details, floating lanterns, front facing, blooming light"),
		),
	)

	if err != nil {
		t.Fatal(err)
	}

	_ = os.WriteFile("result.jpeg", b, os.ModePerm)
	t.Logf("size: %dB", len(b))

}
