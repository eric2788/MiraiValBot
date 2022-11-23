package huggingface

import (
	"os"
	"strings"
	"testing"

	"github.com/eric2788/MiraiValBot/utils/test"
)

const imgRes = `{
    "data": [
        "data:image/png;base64,xxxxxxxxxxx"
    ],
    "flag_index": null,
    "updated_state": null,
    "durations": [
        1.327430248260498
    ],
    "avg_durations": [
        29.61608569352132
    ]
}`

const txtRes = `
{
    "data": [
        "landscape of a cyberpunk abandoned city, destroyed buildings, dystopia, artstation, concept art, painting by greg rutkowski, craig mullins, octane rendering, 8 k, dark atmosphere\n\nlandscape of an asian temple,  planets  around, intricate artwork by Tooth Wu and wlop and beeple. octane render, trending on artstation, greg rutkowski very coherent symmetrical artwork. cinematic, hyper realism, high detail, octane render, 8k\n\nlandscape of a post - apocalyptic dieselpunk city overgrown with lush vegetation, by Luis Royo, by Greg Rutkowski, dark, gritty, intricate, backlit, strong rim light, cover illustration, concept art, volumetric lighting, volumetric atmosphere, sharp focus, octane render, trending on artstation, 8k\n\nlandscape of dark academia aesthetic : : colorful, neon lights, geometric shapes, hard edges, gloomy atmosphere, : a single close up photo - real delicate ceramic black porcelain, high detailed face, photorealism, golden ratio, hyper - realistic 3 d, insanely super detailed, realistic octane render, 1 6 k, minimalistic vulnicura fashion ( by james merry ), jewelry fashion by maiko taked\n"
    ],
    "is_generating": false,
    "duration": 18.62630319595337,
    "average_duration": 19.209084527662707
}`

func init() {
	test.InitTesting()
}

func TestWaifuDiffusier(t *testing.T) {

	if os.Getenv("HUGGING_FACE_TOKEN") == "" {
		t.Skip("no token set")
	}

	api := NewInferenceApi("Nilaier/Waifu-Diffusers",
		Input("group of girls, golden details, gems, fluffy hair, silver hair straight high ponytail+long hair, curly blonde hair, {ginger hair}, purple hair low ponytail curly, clear details, night time, fireworks, casual clothes, raytracing, foreground focus, blurred background, intricate details, floating lanterns, front facing, blooming light"),
	)

	b, err := api.GetResultImage()

	if err != nil {
		t.Fatal(err)
	}

	_ = os.WriteFile("result.jpeg", b, os.ModePerm)
	t.Logf("size: %dB", len(b))

}

func TestBlockedImage(t *testing.T) {
	models := []string{
		"Linaqruf/anything-v3.0",
		"hakurei/waifu-diffusion",
		"Nilaier/Waifu-Diffusers",
	}
	for _, model := range models {
		api := NewInferenceApi(model, Input("1girl, masterpiece, best quality, hyper detailed, Cinematic light, intricate_detail, highres, official art, barefoot, legs, short pants, looking at the viewer"))
		b, err := api.GetResultImage()
		if err != nil {
			t.Log(err)
			continue
		}
		if IsImageBlocked(b) {
			t.Logf("%v image blocked", model)
		} else {
			name := strings.ReplaceAll(model, "/", "-") + "_result.jpeg"
			_ = os.WriteFile(name, b, os.ModePerm)
			t.Logf("%s -> size: %dB", name, len(b))
		}
	}
}

func TestMagicPrompt(t *testing.T) {
	if os.Getenv("HUGGING_FACE_TOKEN") == "" {
		t.Skip("no token set")
	}

	api := NewInferenceApi("Gustavosta/MagicPrompt-Stable-Diffusion", Input("landscape of"))
	txt, err := api.GetGeneratedText()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(txt)
}
