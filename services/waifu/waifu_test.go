package waifu

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/utils/test"
)

func init() {
	test.InitTesting()
}

func TestGetPixivMoe(t *testing.T) {

	file.GenerateConfig()
	file.LoadApplicationYaml()
	Init()

	pixivmoe := &PixelMoe{}
	ids, err := pixivmoe.getPixivIdsByTags([]string{"猫耳", "萝莉"}, 0, 5, false)
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

func TestQueryEncode(t *testing.T) {
	option := NewOptions(
		WithTags("t1", "t2", "t3", "t4"),
		WithKeyword("hawidhaihdi"),
		WithR18(true),
		WithAmount(20),
	)
	r18 := 0
	if option.R18 {
		r18 = 1
	}
	params := &url.Values{
		"tag":     option.Tags,
		"r18":     []string{fmt.Sprint(r18)},
		"num":     []string{fmt.Sprint(option.Amount)},
		"keyword": []string{option.Keyword},
		"size":    []string{"original"},
	}

	t.Log(params.Encode())
}

func TestGetPixivIcon(t *testing.T) {
	url := "https://i.pximg.net/user-profile/img/2022/09/26/02/35/44/23383020_ad04155d3b239285249e6d0837123609_50.jpg"
	b, err := getImageByte(url)
	if err != nil {
		t.Skip(err)
	}
	t.Logf("%d B", len(b))
}

func TestGetDanbooru(t *testing.T) {
	dan := &Danbooru{}
	data, err := dan.GetImages(NewOptions(
		WithTags("loli", "genshin impact"),
		WithAmount(3),
		WithR18(false),
	))

	if err != nil {
		t.Skip(err)
	}

	for _, d := range data {
		t.Logf("title: %s, Tags: %s, R18: %t\n", d.Title, strings.Join(d.Tags, ","), d.R18)
	}

	t.Logf("found %d data", len(data))
}

func TestGetLolicron(t *testing.T) {

	loli := &Lolicron{}

	data, err := loli.GetImages(NewOptions(
		WithTags("萝莉", "兽耳"),
		WithAmount(5),
		WithR18(false),
	))

	if err != nil {
		t.Skip(err)
	}

	for _, d := range data {
		t.Logf("title: %s, Tags: %s, R18: %t\n", d.Title, strings.Join(d.Tags, ","), d.R18)
		if d.R18 {
			t.Fatal("should not have r18")
		}
	}
	t.Logf("found %d data", len(data))
}

func TestGetYuban(t *testing.T) {
	yuban := &Yuban{}

	data, err := yuban.GetImages(NewOptions(
		WithTags("萝莉", "兽耳"),
		WithAmount(5),
		WithR18(false),
	))

	if err != nil {
		t.Skip(err)
	}

	for _, d := range data {
		t.Logf("title: %s, Tags: %s, R18: %t\n", d.Title, strings.Join(d.Tags, ","), d.R18)
		if d.R18 {
			t.Fatal("should not have r18")
		}
	}
	t.Logf("found %d data", len(data))
}

func TestGetAnosuTop(t *testing.T) {
	anosu := &AnosuTop{}
	data, err := anosu.GetImages(NewOptions(
		WithKeyword("genshin"),
		WithAmount(5),
		WithR18(false),
	))

	if err != nil {
		t.Skip(err)
	}

	for _, d := range data {
		t.Logf("title: %s, Tags: %s, R18: %t\n", d.Title, strings.Join(d.Tags, ","), d.R18)
		if d.R18 {
			t.Fatal("should not have r18")
		}
	}
	t.Logf("found %d data", len(data))
}
