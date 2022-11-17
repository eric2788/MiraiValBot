package waifu

import (
	"io"
	"strings"
	"testing"

	"github.com/eric2788/common-utils/request"
)

func TestGetPixivMoe(t *testing.T) {
	pixivmoe := &PixelMoe{}
	ids, err := pixivmoe.getPixivIdsByKeyword("草神", 0, 5, false)
	if err != nil {
		t.Skip(err)
	}
	for _, id := range ids {
		t.Logf("https://pixiv.net/artworks/%d", id)
		data, err := getIllust(id)
		if err != nil {
			t.Log(err)
			continue
		}
		t.Logf("title: %s, tags: %s, url: %s", data.Title, strings.Join(pixivmoe.toArr(data.Tags), ", "), data.Images.Original)
	}
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
		t.Logf("title: %s, tags: %s, url: %s, r18: %t", d.Title, strings.Join(d.Tags, ", "), d.Url, d.R18)
		if d.R18 {
			t.Fatal("should not have r18")
		}
	}
}
