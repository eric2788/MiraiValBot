package urlparser

import (
	"strings"
	"testing"
)

const (
	bvlink    = `https://www.bilibili.com/video/BV1LR4y1y76t/?spm_id_from=333.851.b_7265636f6d6d656e64.5&vd_source=0677b2cd9313952cc0e25879826b251c`
	shortLink = `https://b23.tv/qGyBSoE`
)

func TestParseBV(t *testing.T) {
	matches := biliLinks[0].FindStringSubmatch(bvlink)

	t.Log(strings.Join(matches, ", "))
}

func TestParseShortLink(t *testing.T) {
	s, err := getRedirectLink(shortLink)
	if err != nil {
		t.Skip(err)
	}
	t.Logf("%s => %s", shortLink, s)
}
