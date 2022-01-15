package twitter

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	noTweetLinkPattern = regexp.MustCompile("https:\\/\\/t\\.co\\/\\w+")
)

func TextWithoutTCLink(txt string) string {
	return noTweetLinkPattern.ReplaceAllString(txt, "")
}

func ExtractExtraLinks(data *TweetStreamData) []string {

	extraUrls := make([]string, 0)

	// 分開替代連結和額外連結
	if data.Entities.Urls != nil && len(data.Entities.Urls) > 0 {
		for _, url := range data.Entities.Urls {
			replaced := strings.ReplaceAll(data.Text, url.Url, url.ExpandedUrl)

			// cannot place any urls from data text
			if replaced == data.Text {
				extraUrls = append(extraUrls, url.ExpandedUrl)
			} else {
				data.Text = replaced
			}
		}
	}

	// 取代完畢之後，刪走多餘的 tc link
	data.Text = TextWithoutTCLink(data.Text)

	return extraUrls
}

func GetUserLink(screen string) string {
	return fmt.Sprintf("https://twitter.com/%s", screen)
}

func GetStatusLink(screen, status string) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", screen, status)
}
