package imgtag

import (
	"strings"
	"testing"
)



func TestTagImage(t *testing.T) {
	tags, nsfw, err := GetTagsFromImage("https://preview.redd.it/9gjq7h4szi0a1.jpg?width=640&crop=smart&auto=webp&s=1dcafa0e449331010764731a4f41f095981dad86")
	if err != nil {
		t.Skip(err)
	}
	t.Logf("tags: %s, nsfw: %t", strings.Join(tags, ", "), nsfw)
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