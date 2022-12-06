package urlparser

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/eric2788/common-utils/request"
	"github.com/stretchr/testify/assert"
)

const (
	bvlink    = `https://www.bilibili.com/video/BV1LR4y1y76t/?spm_id_from=333.851.b_7265636f6d6d656e64.5&vd_source=0677b2cd9313952cc0e25879826b251c`
	shortLink = `https://b23.tv/qGyBSoE`
)

func TestParseBV(t *testing.T) {
	matches := parsePattern(bvlink, biliLinks[0])
	assert.Equal(t, "BV1LR4y1y76t", matches[0])
}

func TestParseShortLink(t *testing.T) {
	s, err := getRedirectLink(shortLink)
	if err != nil {
		t.Skip(err)
	}
	t.Logf("%s => %s", shortLink, s)
}

func TestGoQuery(t *testing.T) {
	url := "https://b23.tv/qGyBSoE"
	content, err := request.GetHtml(url)
	if err != nil {
		t.Skipf("解析URL %s 出现错误: %v", url, err)
	}
	docs, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		t.Skipf("解析URL %s 为 html 时出现错误: %v", url, err)
	}
	title := docs.Find("meta[property='og:title']").Text()
	if title == "" {
		title = docs.Find("title").Text()
	}
	thumbnail := docs.Find("meta[property='og:Image']").AttrOr("content", "")
	if thumbnail == "" {
		thumbnail = docs.Find("img").AttrOr("src", "")
	}
	t.Logf("title: %q, thumbnail: %q", title, thumbnail)
}
