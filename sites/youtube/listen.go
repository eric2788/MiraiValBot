package youtube

import (
	"fmt"
	"github.com/eric2788/MiraiValBot/utils/regex"
	"github.com/eric2788/MiraiValBot/utils/request"
	"regexp"
	"strings"
)

var (
	CustomUrlPattern = regexp.MustCompile("(https?:\\/\\/)?(www\\.)?youtube\\.com\\/c\\/(?P<username>[\\w]+)")
	ChannelPattern   = regexp.MustCompile("(https?:\\/\\/)?(www\\.)?youtube\\.com\\/(channel|user)\\/(?P<id>[\\w-]+)")
)

func GetChannelPattern(username string) *regexp.Regexp {
	return regexp.MustCompile("\"browseId\":\"(?P<id>[\\w-]+)\",\"canonicalBaseUrl\":\"\\/c\\/" + username + "\"")
}

// username -> channelID
var channelIdMap = make(map[string]string)

// channelId -> exist
var channelExistMap = make(map[string]bool)

func StartListen(channel string) (bool, error) {
	if !strings.HasPrefix(channel, "UC") {
		return false, fmt.Errorf("不是一個有效的頻道ID")
	}

	if exist, ok := channelExistMap[channel]; ok && !exist {
		return false, fmt.Errorf("頻道不存在")
	} else if !ok {
		// TODO do request check
	}

	// TODO do add and subscribe
	return false, nil
}

func StopListen(channel string) (bool, error) {
	if !strings.HasPrefix(channel, "UC") {
		return false, fmt.Errorf("不是一個有效的頻道ID")
	}

	// TODO do remove and unsubscribe
	return false, nil
}

// GetChannelId get channel id by url
func GetChannelId(url string) (string, error) {

	if ChannelPattern.MatchString(url) {
		return regex.GetParams(ChannelPattern, url)["id"], nil
	}

	if !CustomUrlPattern.MatchString(url) {
		return "", fmt.Errorf("無效的頻道URL")
	}

	username := regex.GetParams(CustomUrlPattern, url)["username"]

	if id, ok := channelIdMap[username]; ok {
		return id, nil
	}

	ex := GetChannelPattern(username)

	body, err := request.GetHtml(fmt.Sprintf("https://www.youtube.com/c/%s", username))

	if err != nil {
		return "", err
	}

	if !ex.MatchString(body) {
		return "", fmt.Errorf("頻道解析失敗")
	}

	channelId := regex.GetParams(ex, body)["id"]
	if !strings.HasPrefix(channelId, "UC") {
		return "", fmt.Errorf("頻道解析失敗, 解析結果為: %s", channelId)
	}

	channelIdMap[username] = channelId
	return channelId, nil
}
