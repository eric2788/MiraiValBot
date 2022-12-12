package responses

import (
	"fmt"
	"math/rand"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/response"
	"github.com/eric2788/MiraiValBot/services/copywriting"
	"github.com/eric2788/MiraiValBot/utils/misc"
)

var tiangouKeywords = []string{
	"田狗",
	"天狗",
	"舔狗",
	"要舔",
	"开舔",
	"开天",
	"好舔",
}

type tiangou struct {
}

func (t *tiangou) ShouldHandle(msg *message.GroupMessage) bool {
	content := msg.ToString()
	return rand.Intn(100)+1 > 65 && misc.ContainsAnyWords(content, tiangouKeywords...)
}

func (t *tiangou) Handle(c *client.QQClient, msg *message.GroupMessage) error {
	list, err := copywriting.GetTianGouList()
	if err != nil {
		return fmt.Errorf("获取天狗列表失败: %v", err)
	}
	random := list[rand.Intn(len(list))]
	return qq.SendGroupMessageByGroup(msg.GroupCode, message.NewSendingMessage().Append(message.NewText(random)))
}

func init() {
	response.AddHandle(&tiangou{})
}
