package chat_reply

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/qq"
)

type randomResponse struct {
}

func (a *randomResponse) Response(msg *message.GroupMessage) (*message.SendingMessage, error) {
	random, err := qq.GetRandomGroupMessage(msg.GroupCode)
	if err != nil {
		logger.Warnf("嘗試發送隨機群消息時出現錯誤: %v, 已略過發送。", err)
		return nil, err
	}
	send := qq.CreateReply(msg)

	for _, ele := range random.Elements {

		switch ele.(type) {
		case *message.ReplyElement:
			continue
		case *message.ForwardElement:
			continue
		default:
			break
		}

		send.Append(ele)
	}
	return send, nil
}
