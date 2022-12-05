package urlparser

import (
	"regexp"
	"strings"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/common-utils/set"
)

const biliVideoInfoURL = "http://api.bilibili.com/x/web-interface/view/detail?bvid=%s"

var (
	biliLinks = []*regexp.Regexp{
		regexp.MustCompile(`https?:\/\/(?:www\.)?bilibili\.com\/video\/(BV\w+)\/?`),
		regexp.MustCompile(`https?:\/\/b23\.tv\/(BV\w+)\/?`),
	}
	shortURLLink = regexp.MustCompile(`https?:\/\/b23\.tv\/(\w+)\/?`)
)

type (
	bilibili struct {
	}

	videoResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		TTL     int    `json:"ttl"`
		Data    struct {
			View struct {
				Bvid        string `json:"bvid"`
				Aid         int64  `json:"aid"`
				TName       string `json:"tname"`
				Title       string `json:"title"`
				Pic         string `json:"pic"`
				PublishDate int64  `json:"pubdate"`
				Ctime       int64  `json:"ctime"`
				Desc        string `json:"desc"`
				Duration    int64  `json:"416"`
				Owner       struct {
					Mid  int64  `json:"mid"`
					Name string `json:"name"`
					Face string `json:"face"`
				} `json:"owner"`
				Stats struct {
				}
			} `json:"View"`
		} `json:"data"`
	}
)

func (b *bilibili) ParseURL(content string) Broadcaster {

	content = b.replaceShortLink(content)

	found := set.NewString()
	for _, patten := range biliLinks {
		matches := parsePattern(content, patten)

		if matches == nil {
			continue
		}

		for _, match := range matches {
			found.Add(match[0])
		}
	}

	if found.Size() == 0 {
		return nil
	}

	return func(bot *client.QQClient, event *message.GroupMessage) {
		for _, bvid := range <-found.Iterator() {
			logger.Debug(bvid)
		}
	}
}

func (b *bilibili) replaceShortLink(content string) string {
	if !shortURLLink.MatchString(content) {
		return content
	}

	for _, matches := range shortURLLink.FindAllStringSubmatch(content, -1) {

		if len(matches) < 2 {
			continue
		}

		if strings.HasPrefix(matches[1], "BV") {
			continue
		}

		link := matches[0]

		s, err := getRedirectLink(link)
		if err != nil {
			logger.Errorf("解析 bilibili 短链接 %s 时出现错误: %v", link, err)
		} else {
			content = strings.ReplaceAll(content, link, s)
		}
	}

	return content
}
