package imgtag

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/eric2788/common-utils/request"
)

const tagURL = "https://nsfwtag.azurewebsites.net/api/tag?limit=0.48&url=%s"

func GetTagsFromImage(imgUrl string) ([]string, bool, error) {
	var dict map[string]float64
	err := request.Get(fmt.Sprintf(tagURL, url.QueryEscape(imgUrl)), &dict)
	if err != nil {
		return nil, false, err
	}
	var tags []string
	nsfw := false
	for key, sample := range dict {
		if key == "rating:safe" {
			nsfw = sample <= 0.55
		} else if key == "rating:questionable" {
			nsfw = sample >= 0.75
		} else if key == "rating:explicit" {
			nsfw = sample >= 0.55
		} else { //filter rating:xxx
			tags = append(tags, strings.ReplaceAll(key, "_", " "))
		}
	}
	return tags, nsfw, nil
}
