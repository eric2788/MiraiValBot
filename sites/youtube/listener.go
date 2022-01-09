package youtube

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/file"
	bc "github.com/eric2788/MiraiValBot/modules/broadcaster"
	"github.com/eric2788/MiraiValBot/utils/regex"
	"github.com/eric2788/MiraiValBot/utils/request"
	"net/http"
	"regexp"
	"strings"
)

var (
	CustomUrlPattern = regexp.MustCompile("(https?:\\/\\/)?(www\\.)?youtube\\.com\\/c\\/(?P<username>[\\w]+)")
	ChannelPattern   = regexp.MustCompile("(https?:\\/\\/)?(www\\.)?youtube\\.com\\/(channel|user)\\/(?P<id>[\\w-]+)")
	listening        = file.DataStorage.Listening
	topic            = func(ch string) string { return fmt.Sprintf("ylive:%s", ch) }
)

func getChannelPattern(username string) *regexp.Regexp {
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
		if res, err := http.Get(fmt.Sprintf("https://youtube.com/channel/%s", channel)); err != nil {
			return false, fmt.Errorf("嘗試檢查頻道時出現錯誤: %v", err)
		} else if res.StatusCode == 404 {
			channelExistMap[channel] = false
			return false, fmt.Errorf("頻道不存在")
		} else if res.StatusCode != 200 {
			return false, fmt.Errorf("嘗試檢查頻道時出現異常: %d, %s", res.StatusCode, res.Status)
		} else { // return 200
			channelExistMap[channel] = true
		}
	}

	file.UpdateStorage(func() {
		listening.Youtube.Add(channel)
	})

	info, err := bot.GetModule(bc.Tag)

	if err != nil {
		return false, err
	}

	broadcaster := info.Instance.(*bc.Broadcaster)

	return broadcaster.Subscribe(topic(channel), MessageHandler)
}

func StopListen(channel string) (bool, error) {
	if !strings.HasPrefix(channel, "UC") {
		return false, fmt.Errorf("不是一個有效的頻道ID")
	}

	if !listening.Youtube.Contains(channel) {
		return false, nil
	}

	file.UpdateStorage(func() {
		listening.Youtube.Delete(channel)
	})

	info, _ := bot.GetModule(bc.Tag)

	broadcaster := info.Instance.(*bc.Broadcaster)

	result := broadcaster.UnSubscribe(topic(channel))

	return result, nil
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

	ex := getChannelPattern(username)

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
