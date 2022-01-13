package twitter

import (
	"fmt"
	"regexp"
)

var (
	noTweetLinkPattern = regexp.MustCompile("https:\\/\\/t\\.co\\/\\w+")
)

func TextWithoutTCLink(txt string) string {
	return noTweetLinkPattern.ReplaceAllString(txt, "")
}

func GetUserLink(screen string) string {
	return fmt.Sprintf("https://twitter.com/%s", screen)
}

func GetStatusLink(screen, status string) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", screen, status)
}
