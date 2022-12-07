package urlparser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Mrs4s/MiraiGo/message"
	bili "github.com/eric2788/MiraiValBot/hooks/sites/bilibili"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/utils/test"
	"github.com/eric2788/common-utils/datetime"
	"github.com/eric2788/common-utils/request"
	"github.com/stretchr/testify/assert"
)

const (
	bvlink    = `https://www.bilibili.com/video/BV1LR4y1y76t/?spm_id_from=333.851.b_7265636f6d6d656e64.5&vd_source=0677b2cd9313952cc0e25879826b251c`
	shortLink = `https://b23.tv/qGyBSoE`
)

var parser = &common{}

func init() {
	test.InitTesting()
}

func TestParseBV(t *testing.T) {
	matches := parsePattern(bvlink, biliLinks[0])
	assert.Equal(t, "BV1LR4y1y76t", matches[0])
}

func TestParseShortLink(t *testing.T) {
	s, err := getRedirectLink(shortLink)
	if err != nil {
		t.Skip(err)
	}
	t.Logf("%s => %s", shortLink, s)
}

func TestBiliParseURL(t *testing.T) {
	url := "https://b23.tv/qGyBSoE"
	data, err := parser.getEmbedData(url)
	if err != nil {
		t.Skip(err)
	}
	t.Logf("%+v", data)
}

func TestYoutubeParseURL(t *testing.T) {
	url := "https://www.youtube.com/watch?v=1G4isv_Fylg"
	data, err := parser.getEmbedData(url)
	if err != nil {
		t.Skip(err)
	}
	t.Logf("%+v", data)
}

func TestChineseText(t *testing.T) {
	txt := "你好世界海外好低啊哇大文件欧服骄傲我就欧服骄傲你发哦我们佛牌我奥法哦呜"
	b := "abcdefg"
	t.Log(len([]rune(txt)), len([]rune(b)))
	assert.Equal(t, 35, len([]rune(txt)))
	assert.Equal(t, 7, len([]rune(b)))
}

func TestBiliParse(t *testing.T) {
	b := &bilibili{}
	url := "https://b23.tv/BV1LR4y1y76t"
	url = b.replaceShortLink(url)

	t.Logf("url is now: %s", url)

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
			t.Skipf("解析bilibili room_id %s 时出现错误: %v", match[0], err)
		} else {
			roomId = id
		}
	}

	t.Logf("found: bvid: %q, roomId: %d", bvid, roomId)

	// 视频解析
	if bvid != "" {
		var resp videoResp
		if err := request.Get(fmt.Sprintf(biliVideoInfoURL, bvid), &resp); err != nil {
			t.Skipf("尝试解析bilibili视频 %s 时出现错误: %v", bvid, err)
		} else if resp.Code != 0 {
			t.Skipf("尝试解析bilibili视频 %s 时出现错误: %s", bvid, resp.Message)
		} else if resp.Data == nil {
			t.Skipf("bilibili视频 %s 的数据为 nil", bvid)
		} else {
			msg := message.NewSendingMessage()
			msg.Append(qq.NewTextfLn("标题: %s", resp.Data.View.Title))
			msg.Append(qq.NewTextfLn("简介: %s", resp.Data.View.Desc))
			msg.Append(qq.NewTextfLn("发布时间: %s", datetime.FormatSeconds(resp.Data.View.PublishDate)))
			msg.Append(qq.NewTextfLn("观看次数: %d | 弹幕数: %d",
				resp.Data.View.Stats.View, resp.Data.View.Stats.Danmaku))
			msg.Append(qq.NewTextfLn("💬: %d | 🔗: %d | 🪙: %d | ⭐: %d",
				resp.Data.View.Stats.Reply, resp.Data.View.Stats.Share,
				resp.Data.View.Stats.Coin, resp.Data.View.Stats.Favourite))

			img, err := test.FakeUploadImageUrl(resp.Data.View.Pic)
			if err != nil {
				logger.Errorf("上传bilibili视频 %s 的封面时出现错误: %v", bvid, err)
			} else {
				msg.Append(img)
			}

			t.Logf("发送消息: \n%s", test.StringifySendingMessage(msg))
		}
	} else if roomId != 0 { // 直播间解析
		info, err := bili.GetRoomInfo(roomId)
		if err != nil {
			t.Skipf("解析 bilibili 直播间 %d 时出现错误: %v", roomId, err)
		} else if info.Code != 0 {
			t.Skipf("解析 bilibili 直播间 %d 时出现错误: %s", roomId, info.Message)
		} else if m, ok := info.Data.(map[string]interface{}); !ok {
			t.Skipf("bilibili 直播间 %d 的数据类型不是 Map", roomId)
		} else {
			data := &bili.RoomInfoData{}
			if err := data.Parse(m); err != nil {
				t.Skipf("解析 bilibili 直播间 %d 数据时出现错误: %v", roomId, err)
			} else {
				msg := message.NewSendingMessage()
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

				img, err := test.FakeUploadImageUrl(data.KeyFrame)
				if err != nil {
					logger.Errorf("为bilibili直播间 %d 上传直播帧图片失败: %v, 将改用直播封面", roomId, err)
					img, err = test.FakeUploadImageUrl(data.UserCover)
					if err != nil {
						logger.Errorf("为bilibili直播间 %d 上传直播封面失败: %v", roomId, err)
					}
				}

				if img != nil {
					msg.Append(img)
				}

				t.Logf("发送消息: \n%s", test.StringifySendingMessage(msg))
			}
		}
	} else {
		t.Skip("没有需要解析的数据")
	}
}
