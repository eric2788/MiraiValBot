package verbose

import (
	"fmt"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/eventhook"
	"github.com/eric2788/MiraiValBot/internal/file"
	qq "github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/internal/redis"
)

type verbose struct {
}

const Tag = "valbot.verbose"

var (
	logger   = utils.GetModuleLogger(Tag)
	instance = &verbose{}
)

func (v *verbose) HookEvent(qqBot *bot.Bot) {

	verboseLiveRoomStatus()

	qqBot.GroupMessageRecalledEvent.Subscribe(func(c *client.QQClient, event *client.GroupMessageRecalledEvent) {

		if !file.DataStorage.Setting.VerboseDelete {
			return
		}

		var who string

		if member := qq.FindGroupMember(event.OperatorUin); member != nil {
			who = member.Nickname
		} else {
			who = fmt.Sprintf("%v", event.OperatorUin)
		}

		msg := message.NewSendingMessage()
		msg.Append(qq.NewTextfLn("%s 所撤回的消息: ", who))
		m, err := qq.GetGroupMessage(event.GroupCode, int64(event.MessageId))
		if err != nil || m == nil {
			msg.Append(qq.NewTextf("获取消息失败: %v", err))
		} else {
			for _, element := range m.Elements {
				msg.Append(element)
			}
		}
		_ = qq.SendGroupMessage(msg)
	})

	qqBot.GroupMessageEvent.Subscribe(func(c *client.QQClient, gm *message.GroupMessage) {

		if !file.DataStorage.Setting.VerboseDelete {
			return
		}

		key := qq.GroupKey(gm.GroupCode, fmt.Sprintf("msg:%d", gm.Id))
		persist := &qq.PersistentGroupMessage{}

		if err := persist.Parse(gm); err != nil {
			logger.Warnf("嘗試序列化群組消息時出現錯誤: %v", err)
			return
		}

		if err := redis.StoreTemp(key, persist); err != nil {
			logger.Warnf("Redis 储存群组消息时出现错误: %v", err)
		} else {
			logger.Infof("Redis 储存临时群组消息成功。")
		}

	})

}

func init() {
	eventhook.RegisterAsModule(instance, "Verbose", Tag, logger)
}
