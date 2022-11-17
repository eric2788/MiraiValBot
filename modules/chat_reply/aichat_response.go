package chat_reply

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
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
		&aichat.QingYunKe{},
		&aichat.TianXing{},
	}

	rand.Seed(time.Now().UnixMicro())
	rand.Shuffle(len(aichats), func(i, j int) { aichats[i], aichats[j] = aichats[j], aichats[i] })

	reply := qq.CreateReply(msg)

	for _, ai := range aichats {

		msg, err := ai.Reply(content)

		if err != nil {
			logger.Errorf("AI %s 回復訊息時出現錯誤: %v, 將使用其他AI", ai.Name(), err)
			continue
		} else {
			a.buildMessage(reply, msg)
			logger.Infof("AI %s 回复信息成功。", ai.Name())
			return reply, nil
		}
	}

	logger.Errorf("所有 AI 均無法回復訊息, 已略過發送。")
	return nil, errors.New("所有 AI 均無法回復訊息。")
}

func (a *AIChatResponse) buildMessage(reply *message.SendingMessage, content string) {

	face, err := regexp.Compile(`\{face:(\d)}`)
	if err != nil {
		logger.Error(err)
		return
	}

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
