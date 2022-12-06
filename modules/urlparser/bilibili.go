package urlparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	bili "github.com/eric2788/MiraiValBot/hooks/sites/bilibili"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/common-utils/datetime"
	"github.com/eric2788/common-utils/request"
)

const biliVideoInfoURL = "http://api.bilibili.com/x/web-interface/view/detail?bvid=%s"

var (
	biliLinks = []*regexp.Regexp{
		regexp.MustCompile(`https?:\/\/(?:\w+\.)?bilibili\.com\/video\/(BV\w+)\/?`),
		regexp.MustCompile(`https?:\/\/b23\.tv\/(BV\w+)\/?`),
	}
	liveLink     = regexp.MustCompile(`https?:\/\/live\.bilibili\.com\/(\d+)\/?`)
	shortURLLink = regexp.MustCompile(`https?:\/\/b23\.tv\/(\w+)\/?`)
)

type (
	bilibili struct {
	}

	videoResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		TTL     int    `json:"ttl"`
		Data    *struct {
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
					View      int64 `json:"view"`
					Danmaku   int64 `json:"danmaku"`
					Reply     int64 `json:"reply"`
					Favourite int64 `json:"favorite"`
					Coin      int64 `json:"coin"`
					Share     int64 `json:"share"`
					Like      int64 `json:"like"`
					DisLike   int64 `json:"dislike"`
				} `json:"stat"`
				Cid   int64 `json:"cid"`
				Pages []struct {
					Part       string `json:"part"`
					FirstFrame string `json:"first_frame"`
					Cid        int64  `json:"cid"`
					Page       int    `json:"page"`
					Vid        string `json:"vid"`
					Weblink    string `json:"weblink"`
				} `json:"pages"`
			} `json:"View"`

			Tags []struct {
				TagId        int    `json:"tag_id"`
				TagName      string `json:"tag_name"`
				Cover        string `json:"cover"`
				HeadCover    string `json:"head_cover"`
				Content      string `json:"content"`
				ShortContent string `json:"short_content"`
				Type         int    `json:"type"`
				Color        string `json:"color"`
			} `json:"Tags"`

			Reply struct {
				Page struct {
					Account int `json:"account"`
					Count   int `json:"count"`
					Num     int `json:"num"`
					Size    int `json:"size"`
				} `json:"page"`

				Replies []struct {
					Type   int   `json:"type"`
					Mid    int64 `json:"mid"`
					RCount int   `json:"rcount"`
					Count  int   `json:"count"`
					Floor  int   `json:"floor"`
					CTime  int64 `json:"ctime"`
					Like   int   `json:"like"`

					Content struct {
						Message string `json:"message"`
						Plat    int    `json:"plat"`
						Devide  string `json:"device"`
					} `json:"content"`
				} `json:"replies"`
			} `json:"Reply"`
		} `json:"data,omitempty"`
	}
)

func (b *bilibili) ParseURL(url string) Broadcaster {

	url = b.replaceShortLink(url)

	bvid, roomId := "", int64(0)
	for _, pattern := range biliLinks {
		matches := parsePattern(url, pattern)
		if matches == nil {
			continue
		}
		bvid = matches[0]
		break
	}

	match := parsePattern(url, liveLink)
	if match != nil {
		if id, err := strconv.ParseInt(match[0], 10, 64); err != nil {
			logger.Errorf("解析bilibili room_id %s 时出现错误: %v", match[0], err)
		} else {
			roomId = id
		}
	}

	if bvid == "" && roomId == 0 {
		logger.Debugf("bilibili 方式无法解析链接: %s, 将使用下一个方式", url)
		return nil
	}

	return func(bot *client.QQClient, event *message.GroupMessage) error {

		// 视频解析
		if bvid != "" {
			var resp videoResp
			if err := request.Get(fmt.Sprintf(biliVideoInfoURL, bvid), &resp); err != nil {
				return fmt.Errorf("尝试解析bilibili视频 %s 时出现错误: %v", bvid, err)
			} else if resp.Code != 0 {
				return fmt.Errorf("尝试解析bilibili视频 %s 时出现错误: %s", bvid, resp.Message)
			} else if resp.Data == nil {
				return fmt.Errorf("bilibili视频 %s 的数据为 nil", bvid)
			} else {
				msg := qq.CreateReply(event)
				msg.Append(qq.NewTextfLn("标题: %s", resp.Data.View.Title))
				msg.Append(qq.NewTextfLn("创作者: %s", resp.Data.View.Owner.Name))

				if len(resp.Data.View.Desc) > 30 {
					resp.Data.View.Desc = resp.Data.View.Desc[:30] + "..."
				}
				msg.Append(qq.NewTextfLn("简介: %s", resp.Data.View.Desc))
				
				msg.Append(qq.NewTextfLn("发布时间: %s", datetime.FormatSeconds(resp.Data.View.PublishDate)))
				msg.Append(qq.NewTextfLn("观看次数: %d | 弹幕数: %d",
					resp.Data.View.Stats.View, resp.Data.View.Stats.Danmaku))
				msg.Append(qq.NewTextfLn("💬: %d | 🔗: %d | 🪙: %d | ⭐: %d",
					resp.Data.View.Stats.Reply, resp.Data.View.Stats.Share,
					resp.Data.View.Stats.Coin, resp.Data.View.Stats.Favourite))

				var tags []string
				for _, tag := range resp.Data.Tags {
					tags = append(tags, tag.TagName)
				}
				msg.Append(qq.NewTextfLn("标签: %s", strings.Join(tags, ", ")))

				img, err := qq.NewImageByUrl(resp.Data.View.Pic)
				if err != nil {
					logger.Errorf("上传bilibili视频 %s 的封面时出现错误: %v", bvid, err)
				} else {
					msg.Append(img)
				}

				return qq.SendGroupMessage(msg)
			}
		} else if roomId != 0 { // 直播间解析
			info, err := bili.GetRoomInfo(roomId)
			if err != nil {
				return fmt.Errorf("解析 bilibili 直播间 %d 时出现错误: %v", roomId, err)
			} else if info.Code != 0 {
				return fmt.Errorf("解析 bilibili 直播间 %d 时出现错误: %s", roomId, info.Message)
			} else if m, ok := info.Data.(map[string]interface{}); !ok {
				return fmt.Errorf("bilibili 直播间 %d 的数据类型不是 Map", roomId)
			} else {
				data := &bili.RoomInfoData{}
				if err := data.Parse(m); err != nil {
					return fmt.Errorf("解析 bilibili 直播间 %d 数据时出现错误: %v", roomId, err)
				} else {
					msg := qq.CreateReply(event)
					msg.Append(qq.NewTextfLn("标题: %s", data.Title))
					status := ""
					switch data.LiveStatus {
					case 0:
						status = "未开播"
					case 1:
						status = "直播中"
					case 2:
						status = "轮播中"
					default:
						status = "未知直播状态: " + fmt.Sprint(data.LiveStatus)
					}
					msg.Append(qq.NewTextfLn("状态: %s", status))
					if data.LiveStatus == 1 {
						msg.Append(qq.NewTextfLn("直播时间: %s", data.LiveTime.Format(datetime.TimeFormat)))
						msg.Append(qq.NewTextfLn("观看人数: %d", data.Online))
					}
					msg.Append(qq.NewTextfLn("分区: %s", data.AreaName))

					img, err := qq.NewImageByUrl(data.KeyFrame)
					if err != nil {
						logger.Errorf("为bilibili直播间 %d 上传直播帧图片失败: %v, 将改用直播封面", roomId, err)
						img, err = qq.NewImageByUrl(data.UserCover)
						if err != nil {
							logger.Errorf("为bilibili直播间 %d 上传直播封面失败: %v", roomId, err)
						}
					}

					if img != nil {
						msg.Append(img)
					}

					return qq.SendGroupMessage(msg)
				}
			}
		} else {
			return fmt.Errorf("没有需要解析的数据")
		}
	}
}

func (b *bilibili) replaceShortLink(url string) string {
	if !shortURLLink.MatchString(url) {
		return url
	}

	matches := shortURLLink.FindStringSubmatch(url)
	if len(matches) < 2 {
		return url
	}

	if strings.HasPrefix(matches[1], "BV") {
		return url
	}

	link := matches[0]

	s, err := getRedirectLink(link)
	if err != nil {
		logger.Errorf("解析 bilibili 短链接 %s 时出现错误: %v", link, err)
	} else {
		url = strings.ReplaceAll(url, link, s)
	}

	return url
}
