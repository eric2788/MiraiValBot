package responses

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/response"
)

type custom struct{}

func (s *custom) ShouldHandle(msg *message.GroupMessage) bool {
	_, ok := file.DataStorage.Responses[msg.ToString()]
	return ok
}

func (s *custom) Handle(c *client.QQClient, msg *message.GroupMessage) {
	res := file.DataStorage.Responses[msg.ToString()]
	m := message.NewSendingMessage().Append(message.NewText(res))
	_ = qq.SendGroupMessageByGroup(msg.GroupCode, m)
}

func init() {
	response.AddHandle(&custom{})
}
