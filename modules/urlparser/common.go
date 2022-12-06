package urlparser

import (
	"encoding/base64"
	"fmt"
	uurl "net/url"
	"strings"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/PuerkitoBio/goquery"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/common-utils/request"
)

type common struct {
}

func (c *common) ParseURL(url string) Broadcaster {
	return func(bot *client.QQClient, event *message.GroupMessage) error {
		content, err := request.GetHtml(url)
		if err != nil {
			return fmt.Errorf("解析URL %s 出现错误: %v", url, err)
		}
		docs, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			return fmt.Errorf("解析URL %s 为 html 时出现错误: %v", url, err)
		}
		title := docs.Find("meta[property='og:title']").AttrOr("content", "")
		if title == "" {
			title = docs.Find("title").Text()
		}

		thumbnail := docs.Find("meta[property='og:Image']").AttrOr("content", "")
		if thumbnail == "" {
			thumbnail = docs.Find("img").AttrOr("src", "")
		}

		if title == "" {
			return fmt.Errorf("无法解析网站 %s 的标题，已略过", url)
		}

		msg := qq.CreateReply(event)
		msg.Append(qq.NewTextfLn("标题: %s", title))
		var img *message.GroupImageElement
		if strings.HasPrefix(thumbnail, "data:image") {
			b64 := strings.Split(thumbnail, "base64,")[1]
			b, err := base64.StdEncoding.DecodeString(b64)
			if err != nil {
				logger.Errorf("URL %s 的封面base64 解析失败: %v, source: %s", url, err, thumbnail)
			}
			img, err = qq.NewImageByByte(b)
			if err != nil {
				logger.Errorf("URL %s 上传封面失败: %v", url, err)
			}
		} else if thumbnail != "" {

			if strings.HasPrefix(thumbnail, "//") {
				// //host.com/static/img/qq.png
				thumbnail = "http:" + thumbnail
			} else if strings.HasPrefix(thumbnail, "/") {
				// /static/img/qq.png
				u, err := uurl.Parse(url)
				if err != nil {
					logger.Errorf("URL %s 解析失败: %v", url, err)
				} else {
					thumbnail = fmt.Sprintf("http://%s%s", u.Host, thumbnail)
				}
			}

			img, err = qq.NewImageByUrl(thumbnail)
			if err != nil {
				logger.Errorf("URL %s 上传封面失败: %v", url, err)
			}
		}

		if img != nil {
			msg.Append(img)
		}

		return qq.SendGroupMessage(msg)
	}
}
