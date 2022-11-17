package waifu

import (
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/utils/test"
	"io"
	"strings"
	"testing"

	"github.com/eric2788/common-utils/request"
)

func TestGetPixivMoe(t *testing.T) {

	test.InitTesting()
	file.GenerateConfig()
	file.LoadApplicationYaml()
	Init()

	pixivmoe := &PixelMoe{}
	ids, err := pixivmoe.getPixivIdsByKeyword("大雄", 0, 5, false)
	if err != nil {
		t.Skip(err)
	}
	for _, id := range ids {
		t.Logf("https://pixiv.net/artworks/%d", id)
		data, err := getIllust(id)
		if err != nil {
			t.Log(err)
			continue
		} else if data == nil || data.Images == nil {
			continue
		}
		t.Logf("title: %s, tags: %s, url: %+v", data.Title, strings.Join(pixivmoe.toArr(data.Tags), ", "), data.Images)
	}
}

func TestGetPixivIcon(t *testing.T) {
	url := "https://i.pximg.net/user-profile/img/2022/09/26/02/35/44/23383020_ad04155d3b239285249e6d0837123609_50.jpg"
	moe := &PixelMoe{}
	b, err := moe.getImageByte(url)
	if err != nil {
		t.Skip(err)
	}
	t.Logf("%d B", len(b))
}

func TestGetLolicron(t *testing.T) {

	loli := &Lolicron{}

	data, err := loli.GetImages(NewOptions(
		WithKeyword("草神"),
		WithAmount(5),
		WithR18(false),
	))

	if err != nil {

		if e, ok := err.(*request.HttpError); ok {
			defer e.Response.Body.Close()
			if b, err := io.ReadAll(e.Response.Body); err == nil {
				t.Log(string(b))
			}
		}

		t.Skip(err)
	}

	for _, d := range data {
		t.Logf("%+v\n", d)
		if d.R18 {
			t.Fatal("should not have r18")
		}
	}
}
