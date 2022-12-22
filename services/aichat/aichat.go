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

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(aichats), func(i, j int) { aichats[i], aichats[j] = aichats[j], aichats[i] })

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
	switch strings.ToLower(mode) {
	case "cn":
		return []AIReply{
			&XiaoAi{},
			&QingYunKe{},
			&TianXing{},
			&MoliYun{},
		}
	case "chatgpt3":
		return []AIReply{
			&Chatgpt3{},
		}
	default:
		return []AIReply{
			&XiaoAi{},
			&QingYunKe{},
			&TianXing{},
			&MoliYun{},
			&Chatgpt3{},
		}
	}
}
