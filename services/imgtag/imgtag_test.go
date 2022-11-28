package imgtag

import (
	"strings"
	"testing"

	"github.com/eric2788/MiraiValBot/internal/file"
)

func TestTagImage(t *testing.T) {
	file.DataStorage.Setting.TagClassifyLimit = 0.5
	tags, err := GetTagsFromImage("https://cdn.discordapp.com/attachments/898123452089778236/1046596783041691700/1280px-.jpg")
	if err != nil {
		t.Skip(err)
	}
	t.Logf("tags: %s", strings.Join(tags, ", "))
	t.Logf("size: %d", len(tags))
}

func TestSearchTags(t *testing.T) {
	tags, err := SearchTags("猫耳")
	if err != nil {
		t.Skip(err)
	}
	for tag, cn := range tags {
		t.Logf("%s: %s", tag, cn)
	}
}
