package huggingface

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/eric2788/common-utils/request"
)

func TestAnythingV3(t *testing.T) {
	api := NewSpaceApi("montagekoko-anything-v3-0",
		"1girl, long hair, car ears")
	bb, err := api.GetResultImages()
	if err != nil {
		t.Skip(err)
	}
	t.Log(len(bb))
}

func TestAnythingV3Img2Img(t *testing.T) {
	t.Skip("too long need 5 minutes")
	img := "https://gchat.qpic.cn/gchatpic_new/0/0-0-42DFE074B3F0C8A416E7D5895AF941D4/0?term=3&is_origin=0&file=42dfe074b3f0c8a416e7d5895af941d4188218-848-1200.jpg"

	b, err := request.GetBytesByUrl(img)
	if err != nil {
		t.Skip(err)
	}

	b64 := fmt.Sprintf("data:%s;base64,", "image/jpeg") + base64.StdEncoding.EncodeToString(b)

	api := NewSpaceApi("akhaliq-anything-v3-0",
		"anything v3",
		"1girl, long hair, car ears",
		7.5,
		25,
		512,
		512,
		0,
		b64,
		0.5,
		BadPrompt,
	)

	bb, err := api.UseWebsocketHandler().GetResultImages()
	if err != nil {
		t.Skip(err)
	}

	t.Log(len(bb))
}

func TestAnythingV3Wss(t *testing.T) {
	t.Skip("too long need 5 minutes")
	api := NewSpaceApi("akhaliq-anything-v3-0",
		"anything v3",
		"1girl, long hair, car ears",
		7.5,
		35,
		720,
		720,
		0,
		nil,
		0.5,
		BadPrompt,
	)

	bb, err := api.UseWebsocketHandler().GetResultImages()
	if err != nil {
		t.Skip(err)
	}

	t.Log(len(bb))
}

func TestRealCugan(t *testing.T) {
	api := NewSpaceApi("shichen1231-Real-CUGAN",
		"data:image/jpeg;base64,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",
		"up2x-latest-denoise2x.pth",
		2,
	)

	bb, err := api.EndPoint("api/predict/").GetResultImages()
	if err != nil {
		t.Skip(err)
	}
	t.Log(len(bb))
}

func TestSessionHashGenerate(t *testing.T) {
	hash := generateSessionHash()
	t.Log(hash)
}
