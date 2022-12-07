package urlparser

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/PuerkitoBio/goquery"
	"github.com/eric2788/MiraiValBot/internal/qq"
)

type (
	common struct {
	}

	embedData struct {
		Host  string
		Title string
		Desc  string
		Thumb string
		Tags  []string
	}
)

func (c *common) ParseURL(url string) Broadcaster {
	return func(bot *client.QQClient, event *message.GroupMessage) error {
		data, err := c.getEmbedData(url)
		if err != nil {
			return err
		}

		title, desc, thumbnail, tags, host := data.Title, data.Desc, data.Thumb, data.Tags, data.Host

		msg := qq.CreateReply(event)
		msg.Append(qq.NewTextfLn("标题: %s", title))

		if desc != "" {
			if len(desc) > 100 && !strings.HasSuffix(desc, "...") {
				desc = desc[:100] + "..."
			}
			msg.Append(qq.NewTextfLn("简介: %s", desc))
		}

		if len(tags) > 0 {
			msg.Append(qq.NewTextfLn("标签: %s", strings.Join(tags, ", ")))
		}

		if thumbnail != "" {

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
			} else {

				if strings.HasPrefix(thumbnail, "//") {
					// //host.com/static/img/qq.png
					thumbnail = "http:" + thumbnail
				} else if strings.HasPrefix(thumbnail, "/") {
					// /static/img/qq.png
					thumbnail = fmt.Sprintf("http://%s%s", host, thumbnail)
				}

				logger.Debugf("即将上传封面: %s", thumbnail)

				img, err = qq.NewImageByUrl(thumbnail)
				if err != nil {
					logger.Errorf("上传封面 %s 失败: %v", thumbnail, err)
				}
			}

			if img != nil {
				msg.Append(img)
			}
		}

		return qq.SendGroupMessage(msg)
	}
}

func (c *common) getEmbedData(url string) (*embedData, error) {
	// same as ParseURL
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求 %s 出现错误: %v", url, err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 %s 出现错误: %v", url, err)
	}
	defer res.Body.Close()
	docs, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("解析URL %s 为 html 时出现错误: %v", url, err)
	}
	title := docs.Find("meta[property='og:title']").AttrOr("content", "")
	if title == "" {
		title = docs.Find("title").Text()
	}
	desc := docs.Find("meta[property='og:description']").AttrOr("content", "")
	if desc == "" {
		desc = docs.Find("meta[name='description']").AttrOr("content", "")
	}

	thumbnail := docs.Find("meta[property='og:image']").AttrOr("content", "")
	if thumbnail == "" {
		thumbnail = docs.Find("img").AttrOr("src", "")
	}

	tags := docs.Find("meta[property='og:video:tag']").
		FilterFunction(func(i int, s *goquery.Selection) bool {
			return s.AttrOr("content", "") != ""
		}).
		Map(func(i int, s *goquery.Selection) string {
			return s.AttrOr("content", "")
		})

	if title == "" {
		return nil, fmt.Errorf("无法解析网站 %s 的标题，已略过", url)
	}

	return &embedData{
		Host:  res.Request.Host,
		Title: title,
		Desc:  desc,
		Tags:  tags,
		Thumb: thumbnail,
	}, nil
}
