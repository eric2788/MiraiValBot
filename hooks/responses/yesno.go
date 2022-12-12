package responses

import (
	"crypto/md5"
	"encoding/binary"
	"math/rand"
	"strings"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/response"
	"github.com/eric2788/MiraiValBot/utils/misc"
)

var (
	questionMarkReplacer = strings.NewReplacer("?", "", "？", "")
)

type yesno struct{}

func (y *yesno) ShouldHandle(msg *message.GroupMessage) bool {
	return misc.YesNoPattern.MatchString(msg.ToString())
}

func (y *yesno) Handle(c *client.QQClient, msg *message.GroupMessage) {
	content := msg.ToString()
	m := message.NewSendingMessage()
	if ans, ok := file.DataStorage.Answers[content]; ok {
		logger.Infof("此问题已被手动设置，因此使用被设置的回答")
		m.Append(message.NewText(y.getResponse(ans)))
	} else {
		ans = y.getQuestionAns(content)
		logger.Infof("自动回答问题 %s 为 %t", content, ans)
		m.Append(message.NewText(y.getResponse(ans)))
	}
	_ = qq.SendGroupMessageByGroup(msg.GroupCode, m)
}

func (y *yesno) getQuestionAns(content string) bool {
	hasher := md5.New()
	question := questionMarkReplacer.Replace(content)
	hashed := hasher.Sum([]byte(question))
	u64 := binary.BigEndian.Uint64(hashed)
	rand.Seed(int64(u64))
	return rand.Intn(2) == 1
}

func (y *yesno) getResponse(is bool) string {
	if is {
		return "确实"
	} else {
		return "并不是"
	}
}

func init() {
	response.AddHandle(&yesno{})
}
