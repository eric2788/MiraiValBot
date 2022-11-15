package chat_reply

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/services/aichat"
)

type AIChatResponse struct {
}

func (a *AIChatResponse) Response(msg *message.GroupMessage) (*message.SendingMessage, error) {

	content := strings.Join(qq.ParseMsgContent(msg.Elements).Texts, "，")

	aichats := []aichat.AIReply{
		&aichat.XiaoAi{},
		&aichat.QingYunKe{},
		&aichat.TianXing{},
	}

	rand.Seed(time.Now().UnixMicro())
	rand.Shuffle(len(aichats), func(i, j int) { aichats[i], aichats[j] = aichats[j], aichats[i] })

	reply := qq.CreateReply(msg)

	for _, aichat := range aichats {

		msg, err := aichat.Reply(content)

		if err != nil {
			logger.Errorf("AI %s 回復訊息時出現錯誤: %v, 將使用其他AI", aichat.Name(), err)
			continue
		} else {
			reply.Append(message.NewText(msg))
			return reply, nil
		}
	}

	logger.Errorf("所有 AI 均無法回復訊息, 已略過發送。")
	return nil, errors.New("所有 AI 均無法回復訊息。")
}
