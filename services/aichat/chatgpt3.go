package aichat

import (
	"fmt"
	"github.com/solywsh/chatgpt"
	"os"
	"time"
)

var ctx *chatgpt.ChatContext

// openAI
type Chatgpt3 struct {
}

func (c *Chatgpt3) Reply(msg string) (string, error) {

	apiKey := os.Getenv("CHATGPT_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("chatgpt api key not set")
	}

	c.ensureInit()

	cli := chatgpt.New(apiKey, "watta-ai", time.Duration(0))
	cli.ChatContext = ctx
	defer cli.Close()

	logger.Debugf("chatgpt3 接收信息: %s", msg)

	resp, err := cli.ChatWithContext(msg)

	if err != nil {
		logger.Debugf("chatgpt3 返回錯誤: %v", err)
		return "", err
	}

	logger.Debugf("chatgpt3 輸出信息: %s", resp)

	return replaces(resp, map[string]string{
		"答：": "",
	}), nil
}

func (c *Chatgpt3) Name() string {
	return "chatgpt3"
}

func (c *Chatgpt3) ensureInit() {
	if ctx == nil {
		ctx = chatgpt.NewContext()
		ctx.SetBackground("群友们都是DD，爱玩抽象和发龙图，还会说怪话，攻击性极强")
	}
}
