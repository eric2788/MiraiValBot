package privatechat

import (
	"strconv"
	"strings"
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/eventhook"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/services/waifu"
	"github.com/eric2788/MiraiValBot/utils/misc"
)

const Tag = "valbot.privatechat"

var (
	logger = utils.GetModuleLogger(Tag)
)

type privateChatResponse struct {
}

func (p *privateChatResponse) HookEvent(bot *bot.Bot) {
	bot.PrivateMessageEvent.Subscribe(func(client *client.QQClient, event *message.PrivateMessage) {

		// 非群友
		if info := qq.FindGroupMember(event.Sender.Uin); info == nil {
			// 無視
			return
		}

		args := strings.Split(event.ToString(), " ")

		if len(args) > 0 && args[0] == "色图" {

			reply := qq.CreatePrivateReply(event)

			amount := 10

			if len(args) > 1 {
				i, err := strconv.Atoi(args[1])
				if err != nil {
					_ = qq.SendPrivateMessage(event.Sender.Uin, reply.Append(qq.NewTextf("无效的数字: %s", args[1])))
					return
				}
				amount = i
			}

			tags := []string{""}

			if len(args) > 2 {
				tags = strings.Split(args[2], ",")
			}

			var search waifu.Searcher
			if len(tags) == 1 {
				search = waifu.WithKeyword(tags[0])
			} else {
				search = waifu.WithTags(tags...)
			}

			imgs, err := waifu.GetRandomImages(
				waifu.NewOptions(
					search,
					waifu.WithAmount(amount),
					waifu.WithR18(true),
				),
			)

			if err != nil {
				_ = qq.SendPrivateMessage(event.Sender.Uin, reply.Append(qq.NewTextf("获取色图失败: %v", err)))
				return
			} else {
				_ = qq.SendPrivateMessage(event.Sender.Uin, reply.Append(qq.NewTextf("正在索取 %s 图片...", strings.Join(tags, ","))))
			}

			forwarder := message.NewForwardMessage()
			wg := &sync.WaitGroup{}

			for _, img := range imgs {
				wg.Add(1)
				go misc.FetchImageToForward(forwarder, img.Url, wg)
			}

			wg.Wait()

			_ = qq.SendPrivateForwardMessage(event.Sender.Uin, forwarder)
		}

	})
}

func init() {
	eventhook.RegisterAsModule(&privateChatResponse{}, "私聊回應", Tag, logger)
}
