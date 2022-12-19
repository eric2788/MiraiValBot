package chat_reply

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/services/aichat"
)

var (
	face   = regexp.MustCompile(`\{face:(\d+)}`)
	AIChat = &AIChatResponse{}
)

type AIChatResponse struct {
}

func (a *AIChatResponse) Response(msg *message.GroupMessage) (*message.SendingMessage, error) {

	content := strings.Join(qq.ParseMsgContent(msg.Elements).Texts, "，")

	res, err := aichat.GetRandomResponse(content)
	if err != nil {
		return nil, err
	}

	reply := qq.CreateReply(msg)
	a.buildMessage(reply, res)

	return reply, nil
}

func (a *AIChatResponse) buildMessage(reply *message.SendingMessage, content string) {

	// ==== add emojis ====
	indexes := face.FindAllStringSubmatchIndex(content, -1)
	if indexes != nil {
		lastTo := 0
		for _, index := range indexes {
			from, to, start, end := index[0], index[1], index[2], index[3]
			if from > 0 {
				reply.Append(message.NewText(content[lastTo:from]))
			}
			faceID, err := strconv.ParseInt(content[start:end], 10, 32)
			if err != nil {
				logger.Errorf("尝试转换表情元素时出现错误: %v", err)
			} else {
				reply.Append(message.NewFace(int32(faceID)))
			}
			lastTo = to
		}

		if lastTo < len(content) {
			reply.Append(message.NewText(content[lastTo:]))
		}
	} else { // normal
		reply.Append(message.NewText(content))
	}
	// ===============
}
