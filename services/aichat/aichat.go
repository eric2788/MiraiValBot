package aichat

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Logiase/MiraiGo-Template/utils"
)

var (
	logger = utils.GetModuleLogger("service.aichat")
)

func GetRandomResponse(content string) (string, error) {

	aichats := GetAIChats(os.Getenv("AI_CHAT_MODE"))

	for _, ai := range aichats {

		msg, err := ai.Reply(content)

		if err != nil {
			logger.Errorf("AI %s 回復訊息時出現錯誤: %v, 將使用其他AI", ai.Name(), err)
			continue
		} else {
			logger.Infof("AI %s 回复信息成功。", ai.Name())
			return msg, nil
		}
	}

	return "", fmt.Errorf("所有 AI 均無法回復訊息。")
}

func GetAIChats(mode string) []AIReply {

	rand.Seed(time.Now().UnixNano())

	switch strings.ToLower(mode) {

	case "random":

		aichats := []AIReply{
			&XiaoAi{},
			&QingYunKe{},
			&TianXing{},
			&MoliYun{},
			&Chatgpt3{},
		}

		rand.Shuffle(len(aichats), func(i, j int) { aichats[i], aichats[j] = aichats[j], aichats[i] })

		return aichats

	case "chatgpt3":
		chats := make([]AIReply, 0)

		aichats := []AIReply{
			&XiaoAi{},
			&QingYunKe{},
			&TianXing{},
			&MoliYun{},
		}

		rand.Shuffle(len(aichats), func(i, j int) { aichats[i], aichats[j] = aichats[j], aichats[i] })

		chats = append(chats, &Chatgpt3{})
		chats = append(chats, aichats...)
		return chats

	default:
		return GetAIChats("random")
	}
}
