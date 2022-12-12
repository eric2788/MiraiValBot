package responses

import (
	"fmt"
	"math/rand"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/chat_reply"
	"github.com/eric2788/MiraiValBot/modules/response"
	"github.com/eric2788/MiraiValBot/utils/misc"
)

type ai struct {
	res *chat_reply.AIChatResponse
}

func (a *ai) ShouldHandle(msg *message.GroupMessage) bool {
	return rand.Intn(50) == 25
}

func (a *ai) Handle(c *client.QQClient, msg *message.GroupMessage) error {
	// 没有文字信息，随机发送龙图?
	if len(qq.ParseMsgContent(msg.Elements).Texts) == 0 {
		send, err := misc.NewRandomDragon()

		if err != nil {
			logger.Errorf("获取龙图失败: %v, 改为发送随机群图片", err)
			send, err = misc.NewRandomImage()
		}

		// 依然失败
		if err != nil {
			return fmt.Errorf("获取图片失败: %v", err)
		}

		return qq.SendGroupMessageByGroup(msg.GroupCode, send)
	}

	// 透过 AI 回复信息
	reply, err := a.res.Response(msg)
	if err != nil {
		return fmt.Errorf("透过 AI 回复对话时出现错误: %v", err)
	} else {

		// create a message with no reply element
		send := message.NewSendingMessage()

		for _, r := range reply.Elements {

			// skip reply and at
			if _, ok := r.(*message.ReplyElement); ok {
				continue
			} else if _, ok = r.(*message.AtElement); ok {
				continue
			}

			send.Append(r)
		}

		return qq.SendGroupMessageByGroup(msg.GroupCode, send)
	}
}

func init() {
	response.AddHandle(&ai{
		res: &chat_reply.AIChatResponse{},
	})
}
